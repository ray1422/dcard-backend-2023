package controller

import (
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ray1422/dcard-backend-2023/controller/pb"
	"github.com/ray1422/dcard-backend-2023/model"
	"github.com/ray1422/dcard-backend-2023/utils"
	"google.golang.org/grpc"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// listHandler handles /list/<id>/<version>?next=[cursor]
func listHandler(c *gin.Context) {
	listIDStr := c.Param("id")
	listID, err := strconv.Atoi(listIDStr)
	if err != nil {
		c.Status(400)
		return
	}
	versionStr := c.Param("version")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.Status(400)
		return
	}
	cursorStr := c.Query("next")
	listNodes, cursor, err := model.GetListNodes(uint(listID), uint(version), cursorStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(404)
			return
		}
		c.Status(500)
		logger.Default.Error(c, "%s", err)

	}
	articles := model.SerializeCursorPagedItems[model.ArticleSerializer](listNodes, cursor)
	c.JSON(200, articles)
}

// ListGRPCServer is the server API for ListGRPCServer service
type ListGRPCServer struct {
	pb.UnimplementedListServiceServer
}

// SetList set list
func (server *ListGRPCServer) SetList(stream pb.ListService_SetListServer) error {
	// TODO
	for u, err := stream.Recv(); ; u, err = stream.Recv() {
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		nodesRPC := u.GetNodes()
		nodes := utils.Map(func(node *pb.Node) model.ListNode {
			// shouldn't happen
			// if node == nil {
			// 	return nil
			// }
			return model.ListNode{
				Version:   u.Version,
				ListID:    u.ListId,
				CreatedAt: time.Now(),
				NodeOrder: node.Order,
				ArticleID: int(node.ArticleId),
			}
		}, nodesRPC)

		err := model.InsertNodes(nodes)
		if err != nil {
			stream.Send(&pb.SetListReply{
				Status: pb.SetListReply_INTERNAL_ERROR,
			})
			log.Println("failed to write list_nodes to db:", err)
		} else {
			err = stream.Send(&pb.SetListReply{
				Status: pb.SetListReply_OK,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// listenGRPC expected to be called by register
func listenGRPC(listener net.Listener) error {
	grpcSrv := grpc.NewServer()
	pb.RegisterListServiceServer(grpcSrv, &ListGRPCServer{})
	err := grpcSrv.Serve(listener)
	return err
}

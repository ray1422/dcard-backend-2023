package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/ray1422/dcard-backend-2023/controller/pb"
	"github.com/ray1422/dcard-backend-2023/model"
	"github.com/ray1422/dcard-backend-2023/utils"
	"github.com/ray1422/dcard-backend-2023/utils/db"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var _ pb.ListServiceServer = ListGRPCServer{}

// ListGRPCServer is the server API for ListGRPCServer service
type ListGRPCServer struct {
	pb.UnimplementedListServiceServer
}

// CreateList gRPC
func (ListGRPCServer) CreateList(ctx context.Context, req *pb.CreateListRequest) (*pb.CreateListReply, error) {
	if req == nil {
		return nil, errors.New("request shouldn't be nil")
	}
	listID, err := model.NewList(req.GetListKey())

	if err != nil {
		log.Println("failed to create list:", err)
		if err == model.ErrDuplicateKeyError {
			return &pb.CreateListReply{
				Status: pb.CreateListReply_DUPLICATE_KEY,
			}, nil
		}
		return &pb.CreateListReply{
			Status: pb.CreateListReply_INTERNAL_ERROR,
		}, nil

	}
	return &pb.CreateListReply{
		ListId: uint32(listID),
		Status: pb.CreateListReply_OK,
	}, nil

}

// DeleteList gRPC
func (ListGRPCServer) DeleteList(ctx context.Context, req *pb.DeleteListRequest) (*pb.DeleteListReply, error) {
	if req == nil {
		return nil, errors.New("request shouldn't be nil")
	}
	listID := req.GetListId()
	err := model.DeleteList(uint(listID))
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.DeleteListReply{
				Status: pb.DeleteListReply_NOT_FOUND,
			}, nil
		}
		return &pb.DeleteListReply{
			Status: pb.DeleteListReply_INTERNAL_ERROR,
		}, nil
	}
	return &pb.DeleteListReply{
		Status: pb.DeleteListReply_OK,
	}, nil
}

// SetList set list gRPC
func (ListGRPCServer) SetList(stream pb.ListService_SetListServer) error {
	for u, err := stream.Recv(); ; u, err = stream.Recv() {
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		nodesRPC := u.GetNodes()
		if nodesRPC == nil {
			return errors.New("request shouldn't be nil")
		}
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

// SetListVersion sets list version. It's expected to be called as all list nodes have been set.
func (ListGRPCServer) SetListVersion(ctx context.Context, req *pb.SetListVersionRequest) (*pb.SetListVersionReply, error) {
	if req == nil {
		return nil, errors.New("request shouldn't be nil")
	}
	err := model.SetListVersion(uint(req.ListId), req.Version)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.SetListVersionReply{
				Status: pb.SetListVersionReply_NOT_FOUND,
			}, nil
		}
		return &pb.SetListVersionReply{
			Status: pb.SetListVersionReply_INTERNAL_ERROR,
		}, nil
	}
	return &pb.SetListVersionReply{
		Status: pb.SetListVersionReply_OK,
	}, nil

}

// listenGRPC expected to be called by register
func listenGRPC(listener net.Listener) error {
	grpcSrv := grpc.NewServer()
	pb.RegisterListServiceServer(grpcSrv, &ListGRPCServer{})
	err := grpcSrv.Serve(listener)
	return err

}

// DeleteListNodeBefore DeleteListNodeBefore
func (ListGRPCServer) DeleteListNodeBefore(ctx context.Context, req *pb.DeleteListNodeBeforeRequest) (*pb.DeleteListNodeBeforeReply, error) {
	if req == nil {
		return nil, errors.New("req shouldn't be nil")
	}
	db := db.GormDB().Where("created_at <= ?", req.Before.AsTime()).Delete(&model.ListNode{})
	err := db.Error
	if err != nil {
		return &pb.DeleteListNodeBeforeReply{Status: pb.DeleteListNodeBeforeReply_INTERNAL_ERROR}, nil
	}
	fmt.Printf("%d nodes has been deleted.\n", db.RowsAffected)
	return &pb.DeleteListNodeBeforeReply{Status: pb.DeleteListNodeBeforeReply_OK}, nil
}

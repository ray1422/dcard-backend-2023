package controller

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/ray1422/dcard-backend-2023/controller/pb"
	"github.com/ray1422/dcard-backend-2023/model"
	"github.com/ray1422/dcard-backend-2023/utils/db"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestCreateAndSetList this function shows the use case of the APIs for the algorithm updating the list.
// First, It creates a list, then adds items to the list, updates the list version, and finally deletes the expired data.

func TestCreateAndSetList(t *testing.T) {
	var myListVersion uint32 = uint32(time.Now().Unix() & 0xffffffff)
	lst := bufconn.Listen(0xfffff + 1)
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lst.Dial()
	}
	go func() {
		err := listenGRPC(lst)
		assert.Nil(t, err)
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	assert.Nil(t, err)
	defer conn.Close()
	client := pb.NewListServiceClient(conn)

	reply, err := client.CreateList(ctx, &pb.CreateListRequest{ListKey: "yet-another-asdf"})
	assert.Nil(t, err)
	fmt.Println(err)
	assert.Equal(t, reply.Status, pb.CreateListReply_OK)
	listID := reply.ListId
	reply, err = client.CreateList(ctx, &pb.CreateListRequest{ListKey: "yet-another-asdf"})
	assert.Equal(t, reply.Status, pb.CreateListReply_DUPLICATE_KEY)

	var cnt int64 = 0
	db.GormDB().Find(&model.List{}).Where("key", "yet-another-asdf").Count(&cnt)
	assert.Equal(t, int64(1), cnt)
	stream, err := client.SetList(ctx)
	assert.Nil(t, err)
	for i := 0; i < 5; i++ {
		err = stream.Send(&pb.SetListRequest{
			ListId:  listID,
			Version: myListVersion,
			Nodes: []*pb.Node{
				{Order: 1, ArticleId: 1},
				{Order: 2, ArticleId: 2},
				{Order: 3, ArticleId: 3},
			},
		})
		assert.Nil(t, err)

		reply, err := stream.Recv()
		if !assert.Nil(t, err) {
			t.Fatalf("%s", err.Error())
		}
		assert.Equal(t, pb.SetListReply_OK, reply.Status)
	}

	// update list version
	retV, err := client.SetListVersion(ctx, &pb.SetListVersionRequest{
		ListId:  listID,
		Version: myListVersion,
	})
	assert.Nil(t, err)
	assert.Equal(t, pb.SetListVersionReply_OK, retV.Status)
	list := model.List{}
	assert.Nil(t, db.GormDB().Take(&list, listID).Error)
	assert.Equal(t, myListVersion, list.Version)

	// There's no foreign key constraint of listID, so it's fine to delete list here.
	delReply, err := client.DeleteList(ctx, &pb.DeleteListRequest{ListId: listID})
	assert.Equal(t, pb.DeleteListReply_OK, delReply.Status)
	delReply, err = client.DeleteList(ctx, &pb.DeleteListRequest{ListId: listID})
	assert.Equal(t, pb.DeleteListReply_NOT_FOUND, delReply.Status)

	// delete list nodes
	ret, err := client.DeleteListNodeBefore(ctx, &pb.DeleteListNodeBeforeRequest{
		Before: timestamppb.Now(),
	})

	assert.Nil(t, err)
	cnt = 0
	db.GormDB().Find(&model.ListNode{}).Count(&cnt)
	assert.Equal(t, pb.DeleteListNodeBeforeReply_OK, ret.Status)
	assert.Zero(t, cnt)
}

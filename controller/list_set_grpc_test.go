package controller

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/ray1422/dcard-backend-2023/controller/pb"
	"github.com/ray1422/dcard-backend-2023/utils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestSetList(t *testing.T) {
	lst, err := net.Listen("tcp", utils.Getenv("GRPC_ADDR", "localhost:3000"))
	go func() {
		err = listenGRPC(lst)
		assert.Nil(t, err)
	}()
	conn, err := grpc.Dial(utils.Getenv("GRPC_ADDR", "localhost:3000"), grpc.WithInsecure())
	assert.Nil(t, err)
	defer conn.Close()
	client := pb.NewListServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := client.SetList(ctx)
	assert.Nil(t, err)
	for i := 0; i < 5; i++ {
		err = stream.Send(&pb.SetListRequest{
			ListId:  10,
			Version: 1,
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
	// TODO chk db write

}

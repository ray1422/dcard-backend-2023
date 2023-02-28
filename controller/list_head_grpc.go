package controller

import (
	"context"
	"errors"

	"github.com/ray1422/dcard-backend-2023/controller/pb"
	"github.com/ray1422/dcard-backend-2023/model"
	"gorm.io/gorm"
)

// Head is actually FindListByKey
func (ListGRPCServer) Head(ctx context.Context, req *pb.HeadRequest) (*pb.HeadReply, error) {
	if req == nil {
		return nil, errors.New("req shouldn't be nil")
	}
	list, err := model.FindListByKey(req.Key)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.HeadReply{Status: pb.HeadReply_NOT_FOUND}, nil
		}
		return &pb.HeadReply{Status: pb.HeadReply_INTERNAL_ERROR}, nil
	}
	return &pb.HeadReply{Status: pb.HeadReply_OK, ListId: uint32(list.ID)}, nil

}

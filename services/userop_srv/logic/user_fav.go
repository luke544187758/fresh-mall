package logic

import (
	"context"
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"luke544187758/userop-srv/dao/mysql"
	"luke544187758/userop-srv/models"
	"luke544187758/userop-srv/proto"
	"time"
)

type UserFavService struct {
}

func NewUserFavService() *UserFavService {
	return &UserFavService{}
}

func (u *UserFavService) GetFavList(ctx context.Context, req *proto.UserFavRequest) (*proto.UserFavListResponse, error) {
	resp := new(proto.UserFavListResponse)
	faves, err := mysql.GetUserFaves(req.UserId)
	if err != nil {
		return nil, err
	}
	for _, item := range faves {
		fav := &proto.UserFavResponse{
			UserId:  item.UserId,
			GoodsId: item.GoodsId,
		}
		resp.Data = append(resp.Data, fav)
	}
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return resp, nil
}

func (u *UserFavService) AddUserFav(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {
	fav := &models.UserFav{
		UserId:  req.UserId,
		GoodsId: req.GoodsId,
		AddTime: sql.NullTime{Time: time.Now()},
	}
	if err := mysql.InsertUserFav(fav); err != nil {
		return &emptypb.Empty{}, err
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return &emptypb.Empty{}, nil
}

func (u *UserFavService) DeleteUserFav(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {

	if err := mysql.DeleteUserFav(req.UserId, req.GoodsId); err != nil {
		return &emptypb.Empty{}, err
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return &emptypb.Empty{}, nil
}

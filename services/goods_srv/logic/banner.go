package logic

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"luke544187758/goods-srv/dao/mysql"
	"luke544187758/goods-srv/models"
	"luke544187758/goods-srv/proto"
	"time"
)

//轮播图
func (g *GoodsService) BannerList(ctx context.Context, req *emptypb.Empty) (*proto.BannerListResponse, error) {
	rsp := &proto.BannerListResponse{}
	list, err := mysql.GetBannerList()
	if err != nil {
		return nil, status.Error(codes.Internal, "query banner list failed")
	}
	if list == nil {
		return nil, status.Error(codes.NotFound, "banner record is not found")
	}

	rsp.Total = int32(len(list))
	for _, v := range list {
		rsp.Data = append(rsp.Data, &proto.BannerResponse{
			Id:    v.ID,
			Index: v.Index,
			Image: v.Image,
			Url:   v.Url,
		})
	}

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return rsp, nil
}
func (g *GoodsService) CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	banner := new(models.Banner)
	banner.Index = req.Index
	banner.Image = req.Image
	banner.Url = req.Url
	banner.IsDeleted = false
	banner.AddTime = time.Now()
	banner.UpdateTime = time.Now()
	lastId, err := mysql.CreateBanner(banner)
	if err != nil {
		return nil, status.Error(codes.Internal, "create new banner failed")
	}
	rsp := &proto.BannerResponse{}
	rsp.Id = int32(lastId)
	rsp.Image = banner.Image
	rsp.Url = banner.Url
	rsp.Index = banner.Index

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return rsp, nil
}
func (g *GoodsService) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	banner, err := mysql.GetBannerWithId(req.Id)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "query the banner record failed")
	}
	if banner == nil {
		return &emptypb.Empty{}, status.Error(codes.NotFound, "the banner record is not found")
	}
	if err := mysql.DeleteBanner(banner.ID); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "delete the banner failed")
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return &emptypb.Empty{}, nil
}
func (g *GoodsService) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	banner, err := mysql.GetBannerWithId(req.Id)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "query the banner record failed")
	}
	if banner == nil {
		return &emptypb.Empty{}, status.Error(codes.NotFound, "the banner record is not found")
	}

	if req.Index != 0 {
		banner.Index = req.Index
	}
	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.Url != "" {
		banner.Url = req.Url
	}
	banner.UpdateTime = time.Now()

	if err := mysql.ModifyBanner(banner); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "modify the banner record failed")
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return &emptypb.Empty{}, nil
}

package logic

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"luke544187758/goods-srv/dao/mysql"
	"luke544187758/goods-srv/models"
	"luke544187758/goods-srv/pkg/snowflake"
	"luke544187758/goods-srv/proto"
	"time"
)

//品牌和轮播图
func (g *GoodsService) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	rsp := &proto.BrandListResponse{}
	list, err := mysql.GetBrandList()
	if err != nil {
		return nil, status.Error(codes.Internal, "get brand list failed")
	}
	if list == nil {
		return nil, status.Error(codes.NotFound, "the brand list is empty")
	}

	rsp.Total = int32(len(list))
	var page int32 = 1
	var perSize int32 = 10
	if req.Page != 0 {
		page = req.Page
	}
	if req.PerSize != 0 {
		perSize = perSize
	}
	if rsp.Total == 0 {
		return rsp, nil
	} else if rsp.Total < page*perSize {
		list = list[(page-1)*perSize : rsp.Total]
	} else {
		list = list[(page-1)*perSize : page*perSize]
	}

	for _, v := range list {
		rsp.Data = append(rsp.Data, &proto.BrandInfoResponse{
			Id:   v.ID,
			Name: v.Name,
			Logo: v.Logo.String,
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
func (g *GoodsService) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	rsp := &proto.BrandInfoResponse{}
	result, err := mysql.GetBrandWithName(req.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, "query brand from db failed")
	}

	if result != nil {
		return nil, status.Error(codes.AlreadyExists, "the brand is already exists")
	}

	brand := new(models.Brand)
	brand.ID = snowflake.GenID()
	brand.Name = req.Name
	brand.Logo.String = req.Logo
	brand.IsDeleted = false
	brand.AddTime = time.Now()
	brand.UpdateTime = time.Now()

	if err := mysql.CreateBrand(brand); err != nil {
		return nil, status.Error(codes.Internal, "create new brand record failed")
	}

	rsp.Id = brand.ID
	rsp.Name = brand.Name
	rsp.Logo = brand.Logo.String

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return rsp, nil
}
func (g *GoodsService) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brand, err := mysql.GetBrand(req.Id)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "query brand from db failed")
	}

	if brand == nil {
		return &emptypb.Empty{}, status.Error(codes.NotFound, "the brand record is not found")
	}

	if err := mysql.DeleteBrand(req.Id); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "delete the brand record from db failed")
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return &emptypb.Empty{}, nil
}
func (g *GoodsService) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brand, err := mysql.GetBrand(req.Id)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "query brand from db failed")
	}

	if brand == nil {
		return &emptypb.Empty{}, status.Error(codes.NotFound, "the brand record is not found")
	}

	if req.Name != "" {
		brand.Name = req.Name
	}

	if req.Logo != "" {
		brand.Logo.String = req.Logo
	}
	brand.UpdateTime = time.Now()

	if err := mysql.ModifyBrand(brand); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "modify the brand record from db failed")
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return &emptypb.Empty{}, nil
}

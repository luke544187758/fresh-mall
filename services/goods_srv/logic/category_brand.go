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

//品牌分类
func (g *GoodsService) CategoryBrandList(ctx context.Context, req *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	rsp := &proto.CategoryBrandListResponse{}
	list, err := mysql.GetCategoryBrandsList()
	if err != nil {
		return nil, status.Error(codes.Internal, "get category brands list failed")
	}
	if list == nil {
		return nil, status.Error(codes.NotFound, "the category brands list is empty")
	}

	rsp.Total = int32(len(list))

	var page int32 = 1
	var pageSize int32 = 10
	if req.Page != 0 {
		page = req.Page
	}
	if req.PerSize != 0 {
		pageSize = req.PerSize
	}
	start := (page - 1) * pageSize

	list = list[start : start+pageSize]
	for _, v := range list {
		sub := &proto.CategoryBrandResponse{}
		sub.Id = v.ID

		brand, _ := mysql.GetBrand(v.BrandId)
		if brand != nil {
			sub.Brand = &proto.BrandInfoResponse{
				Id:   v.BrandId,
				Name: brand.Name,
				Logo: brand.Logo.String,
			}
		}
		category, _ := mysql.GetCategory(v.CategoryId)
		if category != nil {
			sub.Category = &proto.CategoryInfoResponse{
				Id:             v.CategoryId,
				Name:           category.Name,
				ParentCategory: category.ParentCategoryId.Int64,
				Level:          category.Level,
				IsTab:          category.IsTab,
			}
		}
		rsp.Data = append(rsp.Data, sub)
	}
	return rsp, nil
}

//通过category获取brands
func (g *GoodsService) GetCategoryBrandList(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	rsp := &proto.BrandListResponse{}
	brands, err := mysql.GetBrandsWithCategoryId(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "get brands list with category failed")
	}

	if brands == nil {
		return nil, status.Error(codes.NotFound, "the brands list record with category is empty")
	}

	for _, v := range brands {
		rsp.Data = append(rsp.Data, &proto.BrandInfoResponse{
			Id:   v.ID,
			Name: v.Name,
			Logo: v.Logo.String,
		})
	}

	return rsp, nil
}
func (g *GoodsService) CreateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	rsp := &proto.CategoryBrandResponse{}
	category, err := mysql.GetCategory(req.CategoryId)
	if err != nil {
		return nil, status.Error(codes.Internal, "get category record failed")
	}

	if category == nil {
		return nil, status.Error(codes.NotFound, "the category record is not found")
	}

	brand, err := mysql.GetBrand(req.BrandId)
	if err != nil {
		return nil, status.Error(codes.Internal, "get brand record failed")
	}

	if brand == nil {
		return nil, status.Error(codes.NotFound, "the brand record is not found")
	}

	cb := new(models.CategoryBrand)
	cb.ID = snowflake.GenID()
	cb.CategoryId = req.CategoryId
	cb.BrandId = req.BrandId
	cb.IsDeleted = false
	cb.AddTime = time.Now()
	cb.UpdateTime = time.Now()

	if err := mysql.CreateCategoryBrand(cb); err != nil {
		return nil, status.Error(codes.Internal, "create category brand record failed")
	}

	rsp.Id = cb.ID
	rsp.Brand = &proto.BrandInfoResponse{
		Id:   brand.ID,
		Name: brand.Name,
		Logo: brand.Logo.String,
	}
	rsp.Category = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		ParentCategory: category.ParentCategoryId.Int64,
		Level:          category.Level,
		IsTab:          category.IsTab,
	}
	return rsp, nil
}
func (g *GoodsService) DeleteCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	cb, err := mysql.GetCategoryBrandWithId(req.Id)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "query category brand record with id failed")
	}
	if cb == nil {
		return &emptypb.Empty{}, status.Error(codes.NotFound, "the category brand record is not found")
	}

	if err := mysql.DeleteCategoryBrand(req.Id); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "delete category brand record with id failed")
	}

	return &emptypb.Empty{}, nil
}
func (g *GoodsService) UpdateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	cb, err := mysql.GetCategoryBrandWithId(req.Id)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "query category brand record with id failed when modify category brand record")
	}
	if cb == nil {
		return &emptypb.Empty{}, status.Error(codes.NotFound, "the category brand record is not found when modify category brand record")
	}

	category, err := mysql.GetCategory(req.CategoryId)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "get category record failed when modify category brand record")
	}

	if category == nil {
		return &emptypb.Empty{}, status.Error(codes.NotFound, "the category record is not found when modify category brand record")
	}

	brand, err := mysql.GetBrand(req.BrandId)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "get brand record failed when modify category brand record")
	}

	if brand == nil {
		return &emptypb.Empty{}, status.Error(codes.NotFound, "the brand record is not found when modify category brand record")
	}

	cb.CategoryId = req.CategoryId
	cb.BrandId = req.BrandId
	cb.UpdateTime = time.Now()
	if err := mysql.ModifyCategoryBrand(cb); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "modify the category brand record failed")
	}
	return &emptypb.Empty{}, nil
}

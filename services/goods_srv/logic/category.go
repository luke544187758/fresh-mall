package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"luke544187758/goods-srv/dao/mysql"
	"luke544187758/goods-srv/models"
	"luke544187758/goods-srv/pkg/snowflake"
	"luke544187758/goods-srv/proto"
	"time"
)

//商品分类
func (g *GoodsService) GetAllCategoryList(ctx context.Context, req *emptypb.Empty) (*proto.CategoryListResponse, error) {
	rsp := &proto.CategoryListResponse{}
	categories, err := mysql.GetAllCategory()
	if err != nil {
		return nil, status.Error(codes.Internal, "query all category failed")
	}
	if categories == nil {
		return nil, status.Error(codes.NotFound, "record is not found")
	}

	rsp.Total = int32(len(categories))

	var level1 []map[string]interface{}
	var level2 []map[string]interface{}
	var level3 []map[string]interface{}

	for _, v := range categories {
		subRsp := &proto.CategoryInfoResponse{}

		subRsp.Id = v.ID
		subRsp.Name = v.Name
		subRsp.ParentCategory = v.ParentCategoryId.Int64
		subRsp.Level = v.Level
		subRsp.IsTab = v.IsTab

		rsp.Data = append(rsp.Data, subRsp)

		if v.Level == 1 {
			level1 = append(level1, v.CategoryModelToDict())
		} else if v.Level == 2 {
			level2 = append(level2, v.CategoryModelToDict())
		} else if v.Level == 3 {
			level3 = append(level3, v.CategoryModelToDict())
		}
	}
	// 整理数据
	for _, v2 := range level2 {
		var childs []map[string]interface{}
		for _, v3 := range level3 {
			if v3["parent"] == v2["id"] {
				childs = append(childs, v3)
			}
		}
		if childs != nil && len(childs) > 0 {
			v2["sub_category"] = childs
		}
	}

	for _, v1 := range level1 {
		var childs []map[string]interface{}
		for _, v2 := range level2 {
			if v2["parent"] == v1["id"] {
				childs = append(childs, v2)
			}
		}
		if childs != nil && len(childs) > 0 {
			v1["sub_category"] = childs
		}
	}

	bytes, _ := json.Marshal(level1)
	if bytes != nil {
		rsp.JsonData = string(bytes)
	}

	return rsp, nil
}

//获取子分类
func (g *GoodsService) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	rsp := &proto.SubCategoryListResponse{}

	category, err := mysql.GetCategory(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "query category with id failed")
	}
	if category == nil {
		return nil, status.Error(codes.NotFound, "record is not found that query category with id")
	}
	rsp.Info = &proto.CategoryInfoResponse{
		Id:    category.ID,
		Name:  category.Name,
		Level: category.Level,
		IsTab: category.IsTab,
	}
	if category.ParentCategoryId.Valid {
		rsp.Info.ParentCategory = category.ParentCategoryId.Int64
	}

	list, err := mysql.GetSubCategoryList(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "query subcategories failed")
	}
	if list == nil {
		return rsp, nil
	}
	rsp.Total = int32(len(list))
	for _, c := range list {
		crsp := &proto.CategoryInfoResponse{}
		crsp.Id = c.ID
		crsp.Name = c.Name
		if category.ParentCategoryId.Valid {
			crsp.ParentCategory = category.ParentCategoryId.Int64
		}
		crsp.Level = c.Level
		crsp.IsTab = c.IsTab

		rsp.SubCategorys = append(rsp.SubCategorys, crsp)
	}

	return rsp, nil
}
func (g *GoodsService) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	rsp := &proto.CategoryInfoResponse{}

	category := new(models.Category)
	category.ID = snowflake.GenID()
	category.Name = req.Name
	if req.Level != 1 {
		category.ParentCategoryId.Int64 = req.ParentCategory
	}
	category.Level = req.Level
	category.IsTab = req.IsTab
	category.IsDeleted = false
	category.AddTime = time.Now()
	category.UpdateTime = time.Now()

	if err := mysql.CreateCategory(category); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("create new category failed, err:%v", err))
	}

	rsp.Id = category.ID
	rsp.Name = category.Name
	if category.ParentCategoryId.Valid {
		rsp.ParentCategory = category.ParentCategoryId.Int64
	}
	rsp.Level = category.Level
	rsp.IsTab = category.IsTab

	return rsp, nil
}
func (g *GoodsService) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	category, err := mysql.GetCategory(req.Id)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "query category record with id failed")
	}
	if category == nil {
		return &emptypb.Empty{}, status.Error(codes.NotFound, "category record is not found")
	}
	if err := mysql.DeleteCategory(req.Id); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "delete category record with id failed")
	}
	//TODO 根据业务需要，删除响应的category下的商品

	return &emptypb.Empty{}, nil
}
func (g *GoodsService) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {

	category, err := mysql.GetCategory(req.Id)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "query category record with id failed")
	}
	if category == nil {
		return &emptypb.Empty{}, status.Error(codes.NotFound, "category record is not found")
	}
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryId.Int64 = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}

	category.UpdateTime = time.Now()

	if err := mysql.ModifyCategory(category); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "update category detail failed")
	}

	return &emptypb.Empty{}, nil
}

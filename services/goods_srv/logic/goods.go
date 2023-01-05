package logic

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"luke544187758/goods-srv/dao/mysql"
	"luke544187758/goods-srv/models"
	"luke544187758/goods-srv/pkg/snowflake"
	"luke544187758/goods-srv/proto"
	"luke544187758/goods-srv/utils"
	"time"
)

type GoodsService struct {
}

func NewGoodsService() *GoodsService {
	return &GoodsService{}
}

//GoodsList 获取商品列表
func (g *GoodsService) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	rsp := new(proto.GoodsListResponse)
	terms := `WHERE 1=1 `
	var args []interface{}
	if req.KeyWords != "" {
		terms += `AND name LIKE "%?%" `
		args = append(args, req.KeyWords)
	}

	if req.IsHot {
		terms += `AND is_hot = ? `
		args = append(args, req.IsHot)
	}

	if req.IsNew {
		terms += `AND is_new = ? `
		args = append(args, req.IsNew)
	}

	if req.PriceMin != 0 {
		terms += `AND shop_price >= ? `
		args = append(args, req.PriceMin)
	}

	if req.PriceMax != 0 {
		terms += `AND shop_price <= ? `
		args = append(args, req.PriceMax)
	}

	if req.Brand != 0 {
		terms += `AND brand_id = ? `
		args = append(args, req.Brand)
	}

	if req.TopCategory != 0 {
		var ids []int64
		category, err := mysql.GetCategory(req.TopCategory)
		if err == nil && category != nil {
			level := category.Level
			if level == 1 {
				categories, err := mysql.GetCategoriesWithLevel1(req.TopCategory)
				if err == nil && categories != nil {
					for _, v := range categories {
						ids = append(ids, v.ID)
					}
				}
			} else if level == 2 {
				categories, err := mysql.GetCategoriesWithLevel2(req.TopCategory)
				if err == nil && categories != nil {
					for _, v := range categories {
						ids = append(ids, v.ID)
					}
				}
			} else if level == 3 {
				ids = append(ids, req.TopCategory)
			}

			terms += `AND category_id IN (?) `
			args = append(args, ids)
		}
	}

	list, err := mysql.GetGoodsList(terms, args)
	if err != nil {
		return nil, err
	}

	rsp.Total = int32(len(list))

	// 分页
	var page int32 = 1
	var pageSize int32 = 10
	if req.Page != 0 {
		page = req.Page
	}
	if req.PerSize != 0 {
		pageSize = req.PerSize
	}
	offset := pageSize * (page - 1)
	list = list[offset:pageSize]
	for _, v := range list {
		res := proto.GoodsInfoResponse{
			Id:              v.ID,
			CategoryId:      v.CategoryId,
			Name:            v.Name,
			GoodsSn:         v.GoodsSn,
			ClickNum:        v.ClickNum,
			SoldNum:         v.SoldNum,
			FavNum:          v.FavNum,
			MarketPrice:     v.MarketPrice,
			ShopPrice:       v.ShopPrice,
			GoodsBrief:      v.GoodsBrief,
			ShipFree:        v.ShipFree,
			Images:          utils.GetImageUrl(v.Images),
			DescImages:      utils.GetImageUrl(v.DescImages),
			GoodsFrontImage: v.GoodsFrontImage,
			IsNew:           v.IsNew,
			IsHot:           v.IsHot,
			OnSale:          v.OnSale,
			AddTime:         v.AddTime.Unix(),
			Category: &proto.CategoryBriefInfoResponse{
				Id: v.CategoryId,
			},
			Brand: &proto.BrandInfoResponse{
				Id: v.BrandId,
			},
		}
		category, _ := mysql.GetCategory(v.CategoryId)
		if category != nil {
			res.Category.Name = category.Name
		}
		brand, _ := mysql.GetBrand(v.BrandId)
		if brand != nil {
			res.Brand.Name = brand.Name
			res.Brand.Logo = brand.Logo.String
		}
		rsp.Data = append(rsp.Data, &res)
	}

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return rsp, nil
}

//BatchGetGoods 现在用户提交订单有多个商品，你得批量查询商品的信息吧
func (g *GoodsService) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	rsp := &proto.GoodsListResponse{}
	goods, err := mysql.GetBatchGoods(req.Id)
	if err != nil {
		return nil, err
	}
	if goods == nil {
		return nil, nil
	}

	rsp.Total = int32(len(goods))
	for _, v := range goods {
		res := proto.GoodsInfoResponse{
			Id:              v.ID,
			CategoryId:      v.CategoryId,
			Name:            v.Name,
			GoodsSn:         v.GoodsSn,
			ClickNum:        v.ClickNum,
			SoldNum:         v.SoldNum,
			FavNum:          v.FavNum,
			MarketPrice:     v.MarketPrice,
			ShopPrice:       v.ShopPrice,
			GoodsBrief:      v.GoodsBrief,
			ShipFree:        v.ShipFree,
			Images:          utils.GetImageUrl(v.Images),
			DescImages:      utils.GetImageUrl(v.DescImages),
			GoodsFrontImage: v.GoodsFrontImage,
			IsNew:           v.IsNew,
			IsHot:           v.IsHot,
			OnSale:          v.OnSale,
			AddTime:         v.AddTime.Unix(),
			Category: &proto.CategoryBriefInfoResponse{
				Id: v.CategoryId,
			},
			Brand: &proto.BrandInfoResponse{
				Id: v.BrandId,
			},
		}
		category, _ := mysql.GetCategory(v.CategoryId)
		if category != nil {
			res.Category.Name = category.Name
		}
		brand, _ := mysql.GetBrand(v.BrandId)
		if brand != nil {
			res.Brand.Name = brand.Name
			res.Brand.Logo = brand.Logo.String
		}
		rsp.Data = append(rsp.Data, &res)
	}

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return rsp, nil
}
func (g *GoodsService) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	category, err := mysql.GetCategory(req.CategoryId)
	if err != nil {
		return nil, status.Error(codes.Internal, "get product category failed")
	}
	if category == nil {
		return nil, status.Error(codes.NotFound, "product category is not found")
	}

	brand, err := mysql.GetBrand(req.BrandId)
	if err != nil {
		return nil, status.Error(codes.Internal, "get product brand failed")
	}
	if brand == nil {
		return nil, status.Error(codes.NotFound, "product brand is not found")
	}

	goods := &models.Goods{
		ClickNum:        0,
		SoldNum:         0,
		FavNum:          0,
		ID:              snowflake.GenID(),
		CategoryId:      req.CategoryId,
		BrandId:         req.BrandId,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsSn:         req.GoodsSn,
		Name:            req.Name,
		GoodsFrontImage: req.GoodsFrontImage,
		GoodsBrief:      req.GoodsBrief,
		Images:          utils.SetImageUrl(req.Images),
		DescImages:      utils.SetImageUrl(req.DescImages),
		ShipFree:        req.ShipFree,
		OnSale:          req.OnSale,
		IsDeleted:       false,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		AddTime:         time.Now(),
		UpdateTime:      time.Now(),
	}
	if err := mysql.CreateGoods(goods); err != nil {
		return nil, status.Error(codes.Internal, "create product failed")
	}

	rsp := &proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryId,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		AddTime:         goods.AddTime.Unix(),
		Category: &proto.CategoryBriefInfoResponse{
			Id:   category.ID,
			Name: category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo.String,
		},
	}

	//TODO 此处完善库存的设置 - 分布式事务

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return rsp, nil
}
func (g *GoodsService) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	if err := mysql.DeleteGoods(req.Id); err != nil {
		if err == sql.ErrNoRows {
			return &emptypb.Empty{}, status.Error(codes.NotFound, "record is not found")
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return &emptypb.Empty{}, nil
}
func (g *GoodsService) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	category, err := mysql.GetCategory(req.CategoryId)
	if err != nil {
		return nil, status.Error(codes.Internal, "get product category failed")
	}
	if category == nil {
		return nil, status.Error(codes.NotFound, "product category is not found")
	}

	brand, err := mysql.GetBrand(req.BrandId)
	if err != nil {
		return nil, status.Error(codes.Internal, "get product brand failed")
	}
	if brand == nil {
		return nil, status.Error(codes.NotFound, "product brand is not found")
	}

	res, err := mysql.GetGoodsDetail(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "get product detail failed")
	}
	if res == nil {
		return nil, status.Error(codes.NotFound, "product is not found")
	}

	goods := &models.Goods{
		ID:              res.ID,
		CategoryId:      req.CategoryId,
		BrandId:         req.BrandId,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsSn:         req.GoodsSn,
		Name:            req.Name,
		GoodsFrontImage: req.GoodsFrontImage,
		GoodsBrief:      req.GoodsBrief,
		Images:          utils.SetImageUrl(req.Images),
		DescImages:      utils.SetImageUrl(req.DescImages),
		ShipFree:        req.ShipFree,
	}
	if err := mysql.ModifyGoods(goods); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "update product information failed")
	}

	//TODO 此处完善库存的设置 - 分布式事务

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return &emptypb.Empty{}, nil
}
func (g *GoodsService) UpdateGoodsStatus(ctx context.Context, req *proto.GoodsStatusInfo) (*emptypb.Empty, error) {
	res, err := mysql.GetGoodsDetail(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "get product detail failed")
	}
	if res == nil {
		return nil, status.Error(codes.NotFound, "product is not found")
	}

	goodsStatus := &models.GoodsStatus{
		ID:     req.Id,
		IsNew:  req.IsNew,
		IsHot:  req.IsHot,
		OnSale: req.OnSale,
	}

	if err := mysql.ModifyGoodsStatus(goodsStatus); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "update product status failed")
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return &emptypb.Empty{}, nil
}
func (g *GoodsService) GetGoodsDetail(ctx context.Context, req *proto.GoodsInfoRequest) (*proto.GoodsInfoResponse, error) {
	rsp := &proto.GoodsInfoResponse{}
	good, err := mysql.GetGoodsDetail(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("get goods detail failed, err:%v", err))
	}
	if good == nil {
		return nil, status.Error(codes.NotFound, "record is not found")
	}
	// 每次请求商品详情，点击量增加1
	_ = mysql.SetClickNum(good.ClickNum+1, req.Id)
	rsp.Id = good.ID
	rsp.CategoryId = good.CategoryId
	rsp.Name = good.Name
	rsp.GoodsSn = good.GoodsSn
	rsp.ClickNum = good.ClickNum + 1
	rsp.SoldNum = good.SoldNum
	rsp.FavNum = good.FavNum
	rsp.MarketPrice = good.MarketPrice
	rsp.ShopPrice = good.ShopPrice
	rsp.GoodsBrief = good.GoodsBrief
	rsp.ShipFree = good.ShipFree
	rsp.Images = utils.GetImageUrl(good.Images)
	rsp.DescImages = utils.GetImageUrl(good.DescImages)
	rsp.GoodsFrontImage = good.GoodsFrontImage
	rsp.IsNew = good.IsNew
	rsp.IsHot = good.IsHot
	rsp.OnSale = good.OnSale
	rsp.AddTime = good.AddTime.Unix()
	rsp.Category = &proto.CategoryBriefInfoResponse{
		Id: good.CategoryId,
	}
	rsp.Brand = &proto.BrandInfoResponse{
		Id: good.BrandId,
	}

	category, _ := mysql.GetCategory(good.CategoryId)
	if category != nil {
		rsp.Category.Name = category.Name
	}
	brand, _ := mysql.GetBrand(good.BrandId)
	if brand != nil {
		rsp.Brand.Name = brand.Name
		rsp.Brand.Logo = brand.Logo.String
	}

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return rsp, nil
}

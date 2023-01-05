package logic

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"luke544187758/inventory-srv/dao/mysql"
	"luke544187758/inventory-srv/dao/redis"
	"luke544187758/inventory-srv/pkg/snowflake"
	"time"

	"luke544187758/inventory-srv/models"
	"luke544187758/inventory-srv/proto"
)

type InventoryService struct {
}

func NewInventoryService() *InventoryService {
	return &InventoryService{}
}

func (this *InventoryService) SetInventory(ctx context.Context, req *proto.GoodsInventoryInfo) (*emptypb.Empty, error) {
	// 设置库存，如果需要修改库存，也可以使用该接口

	info, err := mysql.GetInventoryWithGoods(req.GoodsId)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "an error occurred in the database operation")
	}
	if info == nil {
		inv := &models.Inventory{
			ID:         snowflake.GenID(),
			Goods:      req.GoodsId,
			Stocks:     req.Num,
			IsDeleted:  false,
			Version:    0,
			AddTime:    time.Now(),
			UpdateTime: time.Now(),
		}
		if err := mysql.CreateInventory(inv); err != nil {
			return &emptypb.Empty{}, status.Error(codes.Internal, "create inventory information failed")
		}
	} else {
		inv := &models.Inventory{
			Goods:      req.GoodsId,
			Stocks:     req.Num,
			UpdateTime: time.Now(),
		}

		if err := mysql.SetInventory(inv); err != nil {
			return &emptypb.Empty{}, status.Error(codes.Internal, "modify inventory information failed")
		}
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return &emptypb.Empty{}, nil
}
func (this *InventoryService) InventoryDetail(ctx context.Context, req *proto.GoodsInventoryInfo) (*proto.GoodsInventoryResponse, error) {

	info, err := mysql.GetInventoryWithGoods(req.GoodsId)
	if err != nil {
		return nil, status.Error(codes.Internal, "an error occurred in the database operation")
	}

	if info == nil {
		return nil, status.Error(codes.NotFound, "inventory record is not found")
	}

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	rsp := &proto.GoodsInventoryResponse{
		GoodsId: info.Goods,
		Num:     info.Stocks,
	}

	return rsp, nil
}
func (this *InventoryService) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	c, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()
	key := fmt.Sprintf("lock:inventory_%d", snowflake.GenID())
	val := snowflake.GenID()
	// 加锁
	if err := redis.LockerAcquire(c, key, val); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "an error occurred when get redis lock")
	}
	var invs []*models.Deduct
	for _, v := range req.GoodsInfo {
		// 查询库存
		info, err := mysql.GetInventoryWithGoods(v.GoodsId)
		if err != nil {
			return &emptypb.Empty{}, status.Error(codes.Internal, "an error occurred in the database operation")
		}

		if info == nil {
			return &emptypb.Empty{}, status.Error(codes.NotFound, "inventory record is not found")
		}
		if info.Stocks < v.Num {
			// 库存不足
			return &emptypb.Empty{}, status.Error(codes.ResourceExhausted, "inventory shortage")
		} else {
			invs = append(invs, &models.Deduct{
				Goods:  v.GoodsId,
				Stocks: info.Stocks - v.Num,
			})
		}
	}

	if err := mysql.SellTransaction(invs); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "an error occurred when modify the inventory record")
	}
	// 释放锁
	if err := redis.LockerRelease(c, key, val); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "an error occurred when release redis lock")
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return &emptypb.Empty{}, nil
}
func (this *InventoryService) ReBack(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	c, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()
	key := fmt.Sprintf("lock:inventory_%d", snowflake.GenID())
	val := snowflake.GenID()
	// 加锁
	if err := redis.LockerAcquire(c, key, val); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "an error occurred when get redis lock")
	}
	var invs []*models.Deduct
	for _, v := range req.GoodsInfo {
		// 查询库存
		info, err := mysql.GetInventoryWithGoods(v.GoodsId)
		if err != nil {
			return &emptypb.Empty{}, status.Error(codes.Internal, "an error occurred in the database operation")
		}

		if info == nil {
			return &emptypb.Empty{}, status.Error(codes.NotFound, "inventory record is not found")
		}
		invs = append(invs, &models.Deduct{
			Goods:  v.GoodsId,
			Stocks: info.Stocks + v.Num,
		})
	}

	if err := mysql.SellTransaction(invs); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "an error occurred when modify the inventory record")
	}
	// 释放锁
	if err := redis.LockerRelease(c, key, val); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "an error occurred when release redis lock")
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return &emptypb.Empty{}, nil
}

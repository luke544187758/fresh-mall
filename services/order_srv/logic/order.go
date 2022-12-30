package logic

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"luke544187758/order-srv/dao/mysql"
	"luke544187758/order-srv/global"
	"luke544187758/order-srv/models"
	"luke544187758/order-srv/pkg/snowflake"
	"luke544187758/order-srv/proto"
	"luke544187758/order-srv/utils"
	"time"
)

type OrderService struct {
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (*OrderService) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	rsp := new(proto.CartItemListResponse)
	list, err := mysql.GetCartItemList(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "get cart information with user id failed")
	}
	if list == nil {
		return nil, status.Error(codes.NotFound, "the cart record with user id is not found")
	}
	rsp.Total = int32(len(list))
	for _, v := range list {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      v.ID,
			UserId:  v.User,
			GoodsId: v.Goods,
			Nums:    v.Nums,
			Checked: v.Checked,
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
func (*OrderService) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {

	item, err := mysql.GetCartItemWithUserId(req.UserId, req.GoodsId)
	if err != nil {
		return nil, status.Error(codes.Internal, "get cart information with user id failed")
	}
	if item != nil { // 存在记录，则更改商品数量
		item.Nums += req.Nums
		if err := mysql.UpdateCartItemCount(item.ID, item.Nums, time.Now()); err != nil {
			return nil, status.Error(codes.Internal, "update cart item count failed")
		}
	} else { // 不存在记录，则创建新的记录
		item = &models.ShoppingCart{
			ID:         snowflake.GenID(),
			User:       req.UserId,
			Goods:      req.GoodsId,
			Nums:       req.Nums,
			Checked:    req.Checked,
			IsDeleted:  false,
			AddTime:    time.Now(),
			UpdateTime: time.Now(),
		}
		if err := mysql.CreateCartItem(item); err != nil {
			return nil, status.Error(codes.Internal, "create new cart item record failed")
		}
	}

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return &proto.ShopCartInfoResponse{
		Id:      item.ID,
		UserId:  item.User,
		GoodsId: item.Goods,
		Nums:    item.Nums,
		Checked: item.Checked,
	}, nil
}
func (*OrderService) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	item, err := mysql.GetCartItemWithId(req.Id)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "get cart information with primary id failed")
	}
	if item == nil {
		return &emptypb.Empty{}, status.Error(codes.NotFound, "the cart record with primary id is not found")
	}

	item.Checked = req.Checked
	if req.Nums > 0 {
		item.Nums = req.Nums
	}
	item.UpdateTime = time.Now()
	if err := mysql.UpdateCartItem(item); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "modify the cart item information failed")
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return &emptypb.Empty{}, nil
}
func (*OrderService) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	item, err := mysql.GetCartItemWithId(req.Id)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "get cart information with primary id failed")
	}
	if item == nil {
		return &emptypb.Empty{}, status.Error(codes.NotFound, "the cart record with primary id is not found")
	}

	if err := mysql.DeleteCartItem(item.ID); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "deleted cart information with primary id failed")
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	return &emptypb.Empty{}, nil
}
func (*OrderService) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	// 获取所有购物车纪录
	cartItems, err := mysql.GetCartItemsWithUserSelect(req.UserId, true)
	if err != nil {
		return nil, status.Error(codes.Internal, "get all selected cart items with user id failed")
	}
	if cartItems == nil {
		return nil, status.Error(codes.NotFound, "get all selected cart items record is not found")
	}
	ids := make([]int64, 0)
	orderGoodsList := make([]*models.OrderGoods, 0)
	goodsSellList := make([]*proto.GoodsInventoryInfo, 0)
	goodsNums := make(map[int64]int32)
	var trans []*models.TransactionParams
	var totalMount float32 = 0

	for _, v := range cartItems {
		ids = append(ids, v.Goods)
		goodsNums[v.Goods] = v.Nums
	}
	// 查询商品的信息
	goodsRes, err := global.GoodsServiceClient.BatchGetGoods(ctx, &proto.BatchGoodsIdInfo{Id: ids})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("goods service is unavailable, err:%v", err))
	}
	for _, v := range goodsRes.Data {
		totalMount += v.ShopPrice * float32(goodsNums[v.Id])
		orderGoodsList = append(orderGoodsList, &models.OrderGoods{
			Goods:      v.Id,
			Nums:       goodsNums[v.Id],
			GoodsPrice: v.ShopPrice,
			GoodsName:  v.Name,
			GoodsImage: v.GoodsFrontImage,
		})
		goodsSellList = append(goodsSellList, &proto.GoodsInventoryInfo{
			GoodsId: v.Id,
			Num:     goodsNums[v.Id],
		})
	}

	// 扣减库存
	_, err = global.InventoryServiceClient.Sell(ctx, &proto.SellInfo{GoodsInfo: goodsSellList})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("inventory service is unavailable, err:%v", err))
	}

	// 创建订单
	order := &models.OrderInfo{
		ID:           snowflake.GenID(),
		User:         req.UserId,
		OrderMount:   totalMount,
		IsDeleted:    false,
		OrderSn:      utils.GenerateOrderSn(req.UserId),
		Address:      req.Address,
		SignerName:   req.Name,
		SignerMobile: req.Mobile,
		Remark:       req.Remark,
		AddTime:      time.Now(),
		UpdateTime:   time.Now(),
	}

	trans = append(trans, mysql.CreateOrderInfoWithTran(order))

	var orderGoodsBatches []interface{}
	// 批量插入订单商品表
	for i := 0; i < len(orderGoodsList); i++ {
		orderGoodsList[i].ID = snowflake.GenID()
		orderGoodsList[i].Order = order.ID
		orderGoodsList[i].IsDeleted = false
		orderGoodsList[i].AddTime = time.Now()
		orderGoodsList[i].UpdateTime = time.Now()
		orderGoodsBatches = append(orderGoodsBatches, orderGoodsList[i])
	}

	trans = append(trans, mysql.BatchCreateOrderGoodsWithTran(orderGoodsBatches))

	// 删除购物车相关的纪录

	trans = append(trans, mysql.DeleteCartItemWithTran(req.UserId))

	// 通过事务执行所有的db操作
	if err = mysql.TransactionExec(trans); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("an error occurred when create new order, err:%v", err))
	}

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return &proto.OrderInfoResponse{
		Id:         order.ID,
		UserId:     order.User,
		OrderSn:    order.OrderSn,
		Remark:     order.Remark,
		OrderMount: totalMount,
	}, nil
}
func (*OrderService) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	rsp := new(proto.OrderListResponse)
	var list []*models.OrderInfo
	var err error
	if req.UserId == 0 {
		list, err = mysql.GetOrderList()
		if err != nil {
			return nil, status.Error(codes.Internal, "get order list failed")
		}
	} else {
		list, err = mysql.GetOrderListWithUserId(req.UserId)
		if err != nil {
			return nil, status.Error(codes.Internal, "get order list with user id failed")
		}
	}

	if list == nil {
		return nil, status.Error(codes.NotFound, "order list record is not found")
	}

	rsp.Total = int32(len(list))
	// 分页
	var page int32 = 1
	var perSize int32 = 10
	if req.Page != 0 {
		page = req.Page
	}
	if req.PerNum != 0 {
		perSize = req.PerNum
	}
	list = list[(page-1)*perSize : page*perSize]
	for _, v := range list {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:         v.ID,
			UserId:     v.User,
			OrderSn:    v.OrderSn,
			PayType:    v.PayType,
			Status:     v.Status,
			Remark:     v.Remark,
			TradeNo:    v.TradeNo,
			OrderMount: v.OrderMount,
			Address:    v.Address,
			Name:       v.SignerName,
			Mobile:     v.SignerMobile,
			PayTime:    v.PayTime.Format("2006-01-02 15:04:05"),
			AddTime:    v.AddTime.Format("2006-01-02 15:04:05"),
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
func (*OrderService) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	rsp := new(proto.OrderInfoDetailResponse)
	item, err := mysql.GetOrderWithId(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "get order item with id failed")
	}
	if item == nil {
		return nil, status.Error(codes.NotFound, "the order item record with id is not found")
	}

	rsp.OrderInfo = &proto.OrderInfoResponse{
		Id:         item.ID,
		UserId:     item.User,
		OrderSn:    item.OrderSn,
		PayType:    item.PayType,
		Status:     item.Status,
		Remark:     item.Remark,
		TradeNo:    item.TradeNo,
		OrderMount: item.OrderMount,
		Address:    item.Address,
		Name:       item.SignerName,
		Mobile:     item.SignerMobile,
		PayTime:    item.PayTime.Format("2006-01-02 15:04:05"),
	}

	list, err := mysql.GetOrderGoodsWithOrderId(item.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "get order goods with order id failed")
	}

	if list == nil {
		return nil, status.Error(codes.NotFound, "the order goods record with order id is not found")
	}

	for _, v := range list {
		rsp.Data = append(rsp.Data, &proto.OrderItemResponse{
			Id:         v.ID,
			OrderId:    v.Order,
			GoodsId:    v.Goods,
			GoodsName:  v.GoodsName,
			GoodsImage: v.GoodsImage,
			GoodsPrice: v.GoodsPrice,
			Nums:       v.Nums,
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
func (*OrderService) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	if err := mysql.UpdateOrderStatus(req.OrderSn, req.Status); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "modify order status with primary id failed")
	}

	if ctx.Err() == context.Canceled {
		return &emptypb.Empty{}, status.Error(codes.Canceled, "request is canceled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return &emptypb.Empty{}, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return &emptypb.Empty{}, nil
}

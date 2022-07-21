package mysql

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"luke544187758/order-srv/models"
	"time"
)

//GetCartItemList 获取用户的所有购物车信息
func GetCartItemList(uid int64) ([]*models.ShoppingCart, error) {
	var result []*models.ShoppingCart
	sqlCmd := `SELECT * FROM shoppingcart WHERE user = ? AND is_deleted = 0`
	if err := db.Select(&result, sqlCmd, uid); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

//GetCartItemWithUserId 根据用户id和产品id查询的商品信息
func GetCartItemWithUserId(uid, gid int64) (*models.ShoppingCart, error) {
	var result models.ShoppingCart
	sqlCmd := `SELECT * FROM shoppingcart WHERE user = ? AND goods = ? AND is_deleted = 0`
	if err := db.Get(&result, sqlCmd, uid, gid); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

//GetCartItemWithId 根据主键id查询的商品信息
func GetCartItemWithId(id int64) (*models.ShoppingCart, error) {
	var result models.ShoppingCart
	sqlCmd := `SELECT * FROM shoppingcart WHERE id = ? AND is_deleted = 0`
	if err := db.Get(&result, sqlCmd, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

//GetCartItemsWithUserSelect 获取用户购物车中选中的商品列表
func GetCartItemsWithUserSelect(uid int64, checked bool) ([]*models.ShoppingCart, error) {
	var result []*models.ShoppingCart
	sqlCmd := `SELECT * FROM shoppingcart WHERE user = ? AND checked = ? AND is_deleted = 0`
	if err := db.Select(&result, sqlCmd, uid, checked); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

//UpdateCartItemCount 更新用户选择商品的数量
func UpdateCartItemCount(id int64, nums int32, updateTime time.Time) error {
	sqlCmd := `UPDATE shoppingcart SET nums = ?,update_time = ? WHERE id = ?`
	_, err := db.Exec(sqlCmd, nums, updateTime, id)
	return err
}

//CreateCartItem 添加商品到购物车
func CreateCartItem(item *models.ShoppingCart) error {
	sqlCmd := `INSERT INTO shoppingcart (id,user,goods,nums,checked,is_deleted,add_time,update_time) VALUES(?,?,?,?,?,?,?,?)`
	_, err := db.Exec(sqlCmd, item.ID, item.User, item.Goods, item.Nums, item.Checked, item.IsDeleted, item.AddTime, item.UpdateTime)
	return err
}

//UpdateCartItem 更新购物车信息
func UpdateCartItem(item *models.ShoppingCart) error {
	sqlCmd := `UPDATE shoppingcart SET nums = ?, checked = ?, update_time = ? WHERE id = ?`
	_, err := db.Exec(sqlCmd, item.Nums, item.Checked, item.UpdateTime, item.ID)
	return err
}

//DeleteCartItem 删除购物车信息
func DeleteCartItem(id int64) error {
	sqlCmd := `UPDATE shoppingcart SET is_deleted = 1 WHERE id = ?`
	_, err := db.Exec(sqlCmd, id)
	return err
}

//DeleteCartItemWithUserId 订单创建后，购物车移除相关的商品 事务操作
func DeleteCartItemWithTran(uid int64) *models.TransactionParams {
	sqlCmd := `UPDATE shoppingcart SET is_deleted = 1 WHERE user = ? AND checked = 1`
	return &models.TransactionParams{
		SqlCmd: sqlCmd,
		Args:   []interface{}{uid},
	}
}

//GetOrderList 获取所有的订单信息
func GetOrderList() ([]*models.OrderInfo, error) {
	var result []*models.OrderInfo
	sqlCmd := `SELECT * FROM orderinfo WHERE is_deleted = 0`
	if err := db.Select(&result, sqlCmd); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

//GetOrderListWithUserId 根据用户id获取所有的订单信息
func GetOrderListWithUserId(uid int64) ([]*models.OrderInfo, error) {
	var result []*models.OrderInfo
	sqlCmd := `SELECT * FROM orderinfo WHERE user = ? AND is_deleted = 0`
	if err := db.Select(&result, sqlCmd, uid); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

//GetOrderWithId 根据主键id获取订单详情
func GetOrderWithId(id int64) (*models.OrderInfo, error) {
	var result models.OrderInfo
	sqlCmd := `SELECT * FROM orderinfo WHERE id = ? AND is_deleted = 0`
	if err := db.Get(&result, sqlCmd, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

//GetOrderGoodsWithOrderId 根据订单号查询订单中的所有商品
func GetOrderGoodsWithOrderId(oid int64) ([]*models.OrderGoods, error) {
	var result []*models.OrderGoods
	sqlCmd := `SELECT * FROM ordergoods WHERE order = ? AND is_deleted = 0`
	if err := db.Select(&result, sqlCmd, oid); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

//BatchCreateOrderGoods 批量插入订单商品
func BatchCreateOrderGoodsWithTran(list []interface{}) *models.TransactionParams {
	sqlCmd := `INSERT INTO ordergoods (id,add_time,is_deleted,update_time,order,goods,goods_name,goods_image,goods_price,nums) VALUES (?),(?),(?)`
	query, args, _ := sqlx.In(sqlCmd, list...)
	return &models.TransactionParams{
		SqlCmd: query,
		Args:   args,
	}
}

//CreateOrderInfo 创建新订单
func CreateOrderInfoWithTran(o *models.OrderInfo) *models.TransactionParams {
	sqlCmd := `INSERT INTO orderinfo (id,user,order_mount,is_deleted,order_sn,address,signer_name,
		signer_mobile,remark,add_time,update_time) VALUES(?,?,?,?,?,?,?,?,?,?,?)`
	args := make([]interface{}, 0)
	args = append(args, o.ID, o.User, o.OrderMount, o.IsDeleted, o.OrderSn, o.Address, o.SignerName,
		o.SignerMobile, o.Remark, o.AddTime, o.UpdateTime)
	return &models.TransactionParams{
		SqlCmd: sqlCmd,
		Args:   args,
	}
}

//UpdateOrderStatus 修改订单状态
func UpdateOrderStatus(id int64, status string) error {
	sqlCmd := `UPDATE orderinfo SET status = ? WHERE id = ?`
	_, err := db.Exec(sqlCmd, status, id)
	return err
}

//TransactionExec 事务执行所有的语句
func TransactionExec(params []*models.TransactionParams) error {
	tx, err := db.Beginx() // 开启事务
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()

	for _, param := range params {
		_, err = tx.Exec(param.SqlCmd, param.Args...)
		if err != nil {
			return err
		}
	}
	return nil
}

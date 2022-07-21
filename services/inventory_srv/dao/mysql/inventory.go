package mysql

import (
	"database/sql"
	"luke544187758/inventory-srv/models"
)

//SetInventory 修改库存
func SetInventory(inv *models.Inventory) error {
	sqlCmd := `UPDATE inventory SET stocks = ?, update_time = ? WHERE goods = ?`
	_, err := db.Exec(sqlCmd, inv.Stocks, inv.UpdateTime, inv.Goods)
	return err
}

//GetInventoryWithGoods 查询商品的库存
func GetInventoryWithGoods(goods int64) (*models.Inventory, error) {
	var inv models.Inventory
	sqlCmd := `SELECT * FROM inventory WHERE goods = ?`
	if err := db.Get(&inv, sqlCmd, goods); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &inv, nil
}

//CreateInventory 创建库存信息
func CreateInventory(inv *models.Inventory) error {
	sqlCmd := `INSERT INTO inventory (goods,stocks,is_deleted,version,add_time,update_time) VALUES (?,?,?,?,?,?)`
	_, err := db.Exec(sqlCmd, inv.Goods, inv.Stocks, inv.IsDeleted, inv.Version, inv.AddTime, inv.UpdateTime)
	return err
}

func SellTransaction(invs []*models.Deduct) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	for _, v := range invs {
		sqlCmd := `UPDATE inventory SET stocks = ? WHERE goods = ?`
		_, err = tx.Exec(sqlCmd, v.Stocks, v.Goods)
		if err != nil {
			return err
		}
	}
	return nil
}

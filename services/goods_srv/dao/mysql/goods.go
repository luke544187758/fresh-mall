package mysql

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"luke544187758/goods-srv/models"
)

//GetGoodsList 查询商品列表，按照指定条件
func GetGoodsList(terms string, args []interface{}) ([]*models.Goods, error) {
	var goods []*models.Goods
	sqlCmd := fmt.Sprintf("select * from goods %s", terms)
	err := db.Select(&goods, sqlCmd, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return goods, nil
}

//GetBatchGoods 批量获取商品详情
func GetBatchGoods(ids []int64) ([]*models.Goods, error) {
	var goods []*models.Goods
	query, args, err := sqlx.In("SELECT * FROM goods WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	if err = db.Select(&goods, query, args); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return goods, nil
}

//DeleteGoods 根据id删除商品
func DeleteGoods(id int64) error {
	sqlCmd := `UPDATE goods SET is_deleted = 1 WHERE id = ?`
	_, err := db.Exec(sqlCmd, id)
	return err
}

//GetGoodsDetail 获取商品详情
func GetGoodsDetail(id int64) (*models.Goods, error) {
	var goods models.Goods
	sqlCmd := `SELECT * FROM goods WHERE id = ?`
	if err := db.Get(&goods, sqlCmd, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &goods, nil
}

//SetClickNum 每次请求商品详情，更新点击量
func SetClickNum(clickNum int32, id int64) error {
	sqlCmd := `UPDATE goods SET click_num = ? WHERE id = ?`
	_, err := db.Exec(sqlCmd, clickNum, id)
	return err
}

//CreateGoods 添加新商品
func CreateGoods(g *models.Goods) error {
	sqlCmd := `INSERT INTO goods (id,add_time,is_deleted,update_time,category_id,brand_id,
			on_sale,goods_sn,name,click_num,sold_num,fav_num,market_price,shop_price,
			goods_brief,ship_free,images,desc_images,goods_front_image,is_new,is_hot) 
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := db.Exec(sqlCmd, g.ID, g.AddTime, g.IsDeleted, g.UpdateTime, g.CategoryId, g.BrandId,
		g.OnSale, g.GoodsSn, g.Name, g.ClickNum, g.SoldNum, g.FavNum, g.MarketPrice, g.ShopPrice,
		g.GoodsBrief, g.ShipFree, string(g.Images), string(g.DescImages), g.GoodsFrontImage, g.IsNew, g.IsHot)
	return err
}

func ModifyGoods(g *models.Goods) error {
	sqlCmd := `UPDATE goods SET category_id=?, brand_id=?, 
                 goods_sn=?, name=?, market_price=?,shop_price=?, goods_brief=?, ship_free=?, 
                 images=?, desc_images=?, goods_front_image=? WHERE id=?`
	_, err := db.Exec(sqlCmd, g.CategoryId, g.BrandId, g.GoodsSn, g.Name, g.MarketPrice,
		g.ShopPrice, g.GoodsBrief, g.ShipFree, string(g.Images), string(g.DescImages), g.GoodsFrontImage)
	return err
}

func ModifyGoodsStatus(s *models.GoodsStatus) error {
	sqlCmd := `UPDATE goods SET is_hot=?, is_new=?, on_sale=? WHERE id = ?`
	_, err := db.Exec(sqlCmd, s.IsHot, s.IsNew, s.OnSale, s.ID)
	return err
}

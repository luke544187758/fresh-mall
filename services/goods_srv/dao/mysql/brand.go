package mysql

import (
	"database/sql"
	"luke544187758/goods-srv/models"
)

//GetBrand 根据id查询对应的品牌
func GetBrand(bid int64) (*models.Brand, error) {
	brand := new(models.Brand)
	sqlCmd := `SELECT id,name,logo,is_deleted,add_time,update_time FROM brands WHERE id = ?`
	err := db.Get(brand, sqlCmd, bid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return brand, nil
}

//GetBrandList 获取品牌列表
func GetBrandList() ([]*models.Brand, error) {
	var list []*models.Brand
	sqlCmd := `SELECT * FROM brands`
	if err := db.Select(&list, sqlCmd); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return list, nil
}

//CreateBrand 添加新的品牌
func CreateBrand(b *models.Brand) error {
	sqlCmd := `INSERT INTO brands (id, name, logo, is_deleted, add_time, update_time) VALUES (?,?,?,?,?,?)`
	_, err := db.Exec(sqlCmd, b.ID, b.Name, b.Logo, b.IsDeleted, b.AddTime, b.UpdateTime)
	return err
}

//GetBrandWithName 根据name查找对应的品牌
func GetBrandWithName(name string) (*models.Brand, error) {
	var brand models.Brand
	sqlCmd := `SELECT * FROM brands WHERE name = ?`
	if err := db.Get(&brand, sqlCmd, name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &brand, nil
}

//DeleteBrand 根据id删除对应的品牌
func DeleteBrand(id int64) error {
	sqlCmd := `UPDATE brands SET is_deleted = 1 WHERE id = ?`
	_, err := db.Exec(sqlCmd, id)
	return err
}

//ModifyBrand 修改品牌信息
func ModifyBrand(b *models.Brand) error {
	sqlCmd := `UPDATE brands SET name = ?, logo = ?, update_time = ? WHERE id = ?`
	_, err := db.Exec(sqlCmd, b.Name, b.Logo, b.UpdateTime, b.ID)
	return err
}

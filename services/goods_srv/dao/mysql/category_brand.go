package mysql

import (
	"database/sql"
	"luke544187758/goods-srv/models"
)

//GetCategoryBrandsList 获取品牌分类列表
func GetCategoryBrandsList() ([]*models.CategoryBrand, error) {
	var list []*models.CategoryBrand
	sqlCmd := `SELECT * FROM goodscategorybrand`
	if err := db.Select(&list, sqlCmd); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return list, nil
}

//GetBrandsWithCategoryId 根据分类查找对应的品牌列表
func GetBrandsWithCategoryId(cid int64) ([]*models.Brand, error) {
	var list []*models.Brand
	sqlCmd := `SELECT b.id,b.name,b.logo,b.is_deleted,b.add_time,b.update_time FROM goodscategorybrand g 
    			LEFT JOIN brands b ON g.brand_id = b.id WHERE category_id = ?`
	err := db.Select(&list, sqlCmd)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return list, nil
}

//CreateCategoryBrand 添加品牌分类记录
func CreateCategoryBrand(cb *models.CategoryBrand) error {
	sqlCmd := `INSERT INTO goodscategorybrand (id,category_id,brand_id,is_deleted,add_time,update_time) VALUES (?,?,?,?,?,?)`
	_, err := db.Exec(sqlCmd, cb.ID, cb.CategoryId, cb.BrandId, cb.IsDeleted, cb.AddTime, cb.UpdateTime)
	return err
}

//DeleteCategoryBrand 根据id删除对应品牌分类
func DeleteCategoryBrand(id int64) error {
	sqlCmd := `UPDATE goodscategorybrand SET is_deleted = 1 WHERE id = ?`
	_, err := db.Exec(sqlCmd, id)
	return err
}

//GetCategoryBrandWithId 根据id查找对应的品牌分类
func GetCategoryBrandWithId(id int64) (*models.CategoryBrand, error) {
	var cb models.CategoryBrand
	sqlCmd := `SELECT * FROM goodscategorybrand WHERE id = ?`
	if err := db.Get(&cb, sqlCmd, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &cb, nil
}

//ModifyCategoryBrand 更改品牌分类信息
func ModifyCategoryBrand(cb *models.CategoryBrand) error {
	sqlCmd := `UPDATE goodscategorybrand SET category_id = ?, brand_id = ?, update_time = ? WHERE id = ?`
	_, err := db.Exec(sqlCmd, cb.CategoryId, cb.BrandId, cb.UpdateTime, cb.ID)
	return err
}

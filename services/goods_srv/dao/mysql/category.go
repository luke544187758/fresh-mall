package mysql

import (
	"database/sql"
	"luke544187758/goods-srv/models"
)

//GetCategory 根据id查询对应的分类
func GetCategory(cid int64) (*models.Category, error) {
	category := new(models.Category)
	sqlCmd := `SELECT id,name,parent_category_id,level,is_deleted,is_tab,add_time,update_time FROM category WHERE id = ?`
	err := db.Get(category, sqlCmd, cid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return category, nil
}

func GetSubCategoryList(cid int64) ([]*models.Category, error) {
	var list []*models.Category
	sqlCmd := `SELECT * FROM category WHERE parent_category_id = ?`
	if err := db.Select(&list, sqlCmd, cid); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return list, nil

}

//GetCategoriesWithLevel1 查询一级分类下的所有类别
func GetCategoriesWithLevel1(pid int64) ([]*models.Category, error) {
	var categories []*models.Category
	sqlCmd := `SELECT * FROM category WHERE parent_category_id IN (SELECT id FROM category WHERE parent_category_id = ?)`
	err := db.Select(&categories, sqlCmd, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return categories, nil
}

//GetCategoriesWithLevel2 查询二级分类下的所有类别
func GetCategoriesWithLevel2(pid int64) ([]*models.Category, error) {
	var categories []*models.Category
	sqlCmd := `SELECT * FROM category WHERE parent_category_id = ?`
	err := db.Select(&categories, sqlCmd, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return categories, nil
}

//GetAllCategory 获取所有的分类
func GetAllCategory() ([]*models.Category, error) {
	var categories []*models.Category
	sqlCmd := `SELECT * FROM category`
	err := db.Select(&categories, sqlCmd)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return categories, nil
}

//CreateCategory 创建新分类
func CreateCategory(c *models.Category) error {
	sqlCmd := `INSERT INTO category (id, name, parent_category, level, is_tab, is_deleted, add_time, update_time)
			VALUES (?,?,?,?,?,?,?,?)`
	_, err := db.Exec(sqlCmd, c.ID, c.Name, c.ParentCategoryId, c.Level, c.IsTab, c.IsDeleted, c.AddTime, c.UpdateTime)
	return err
}

//DeleteCategory 根据id删除对应的分类记录
func DeleteCategory(id int64) error {
	sqlCmd := `UPDATE category SET is_deleted = 1 WHERE id = ?`
	_, err := db.Exec(sqlCmd, id)
	return err
}

//ModifyCategory 更新分类
func ModifyCategory(c *models.Category) error {
	sqlCmd := `UPDATE category SET name=?, parent_category=?, level=?, update_time = ? WHERE id =?`
	_, err := db.Exec(sqlCmd, c.Name, c.ParentCategoryId, c.Level, c.UpdateTime, c.ID)
	return err
}

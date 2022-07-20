package mysql

import (
	"database/sql"
	"luke544187758/goods-srv/models"
)

//GetBannerList 获取轮播图列表
func GetBannerList() ([]*models.Banner, error) {
	var banners []*models.Banner
	sqlCmd := `SELECT * FROM banner`
	if err := db.Select(&banners, sqlCmd); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return banners, nil
}

//CreateBanner 创建新的轮播图
func CreateBanner(b *models.Banner) (int64, error) {
	sqlCmd := `INSERT INTO banner (banner.index, image, url, is_deleted, add_time, update_time) VALUES (?,?,?,?,?,?)`
	res, err := db.Exec(sqlCmd, b.Index, b.Image, b.Url, b.IsDeleted, b.AddTime, b.UpdateTime)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

//DeleteBanner 根据id删除对应的轮播图
func DeleteBanner(id int32) error {
	sqlCmd := `UPDATE banner SET is_deleted = 1 WHERE id = ?`
	_, err := db.Exec(sqlCmd, id)
	return err
}

//GetBannerWithId 根据id查找对应的轮播图
func GetBannerWithId(bid int32) (*models.Banner, error) {
	var banner models.Banner
	sqlCmd := `SELECT * FROM banner WHERE id = ?`
	if err := db.Get(&banner, sqlCmd, bid); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &banner, nil
}

//ModifyBanner 修改轮播图信息
func ModifyBanner(b *models.Banner) error {
	sqlCmd := `UPDATE banner SET image = ?, url = ?, banner.index = ?, update_time = ? WHERE id = ?`
	_, err := db.Exec(sqlCmd, b.Image, b.Url, b.Index, b.UpdateTime, b.ID)
	return err
}

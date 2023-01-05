package mysql

import (
	"database/sql"
	"luke544187758/userop-srv/models"
)

func GetUserFaves(uid int64) (faves []*models.UserFav, err error) {
	sqlCmd := `SELECT user_id,goods_id FROM userfav WHERE user_id = ?`
	if err = db.Select(&faves, sqlCmd, uid); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return nil, err
	}
	return faves, nil
}

func InsertUserFav(fav *models.UserFav) error {
	sqlCmd := `INSERT INTO userfav (user_id,goods_id,add_time) VALUES (?,?,?,?,?)`
	_, err := db.Exec(sqlCmd, fav.UserId, fav.GoodsId, fav.AddTime)
	return err
}

func DeleteUserFav(uid, goodsId int64) error {
	sqlCmd := `DELETE userfav WHERE user_id = ? AND goods_id = ?`
	_, err := db.Exec(sqlCmd, uid, goodsId)
	return err
}

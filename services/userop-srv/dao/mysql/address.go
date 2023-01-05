package mysql

import (
	"database/sql"
	"go.uber.org/zap"
	"luke544187758/userop-srv/models"
)

func GetAddresses(uid int64) (res []*models.Address, err error) {
	sqlCmd := `SELECT * FROM address WHERE user_id = ?`
	if err = db.Select(&res, sqlCmd, uid); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no user in db")
			err = nil
		}
		return nil, err
	}
	return res, nil
}

func InsertAddress(addr *models.Address) error {
	sqlCmd := `INSERT INTO address (id,user_id,province,city,district,address,signer_name,signer_mobile,add_time,update_time,is_deleted) VALUES (?,?,?,?,?,?,?,?,?,?,?)`
	_, err := db.Exec(sqlCmd, addr.Id, addr.UserId, addr.Province, addr.City, addr.District, addr.Address, addr.SignerName, addr.SignerMobile, addr.AddTime, addr.UpdateTime, addr.IsDeleted)
	return err
}

func UpdateAddress(addr *models.Address) error {
	sqlCmd := `UPDATE address SET province = ?,city = ?, district = ?, address = ?, signer_name = ?, signer_mobile = ?, update_time = ?`
	_, err := db.Exec(sqlCmd, addr.Province, addr.City, addr.District, addr.Address, addr.SignerName, addr.SignerMobile, addr.UpdateTime)
	return err
}

func DeleteAddress(id int64) error {
	sqlCmd := `UPDATE address SET is_deleted = ? WHERE id = ?`
	_, err := db.Exec(sqlCmd, true, id)
	return err
}

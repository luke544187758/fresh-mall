package mysql

import (
	"database/sql"
	"luke544187758/userop-srv/models"
)

func GetMessages(uid int64) (messages []*models.Message, err error) {
	sqlCmd := `SELECT id,user_id,message_type,subject,message,add_time,update_time,is_deleted FROM leavingmessages WHERE user_id = ?`
	if err = db.Select(&messages, sqlCmd, uid); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return nil, err
	}
	return messages, nil
}

func InsertMessage(msg *models.Message) (err error) {
	sqlCmd := `INSERT INTO leavingmessages (id,user_id,message_type,subject,message,add_time,update_time,is_deleted) VALUES (?,?,?,?,?,?,?,?)`
	_, err = db.Exec(sqlCmd, msg.Id, msg.UserId, msg.MessageType, msg.Subject, msg.Message, msg.AddTime, msg.UpdateTime, msg.IsDeleted)
	return err
}

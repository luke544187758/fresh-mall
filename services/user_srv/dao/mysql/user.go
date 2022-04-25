package mysql

import (
	"database/sql"
	"go.uber.org/zap"
	"luke544187758/user-srv/models"
	"luke544187758/user-srv/utils"
)

//GetUserList 获取用户列表
func GetUserList(page uint32, pageSize uint32) (users []*models.User, err error) {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	sqlCmd := `select id,nick_name,password,mobile,gender,address,birthday,role from user limit ?,?`
	if err = db.Select(&users, sqlCmd, (page-1)*pageSize, pageSize); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no user in db")
			err = nil
		}
		return nil, err
	}
	return users, nil
}

//GetUserByMobile 根据手机号获取用户信息
func GetUserByMobile(mobile string) (*models.User, error) {
	post := new(models.User)
	sqlCmd := "select id,nick_name,password,mobile,gender,address,birthday,role from user where mobile = ?"
	err := db.Get(post, sqlCmd, mobile)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return post, err
}

//GetUserByID 根据id获取用户信息
func GetUserByID(id int64) (*models.User, error) {
	post := new(models.User)
	sqlCmd := "select id,nick_name,password,mobile,gender,address,birthday,role from user where id = ?"
	err := db.Get(post, sqlCmd, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return post, err
}

//CreateUser 创建用户
func CreateUser(nickname, mobile, password string) (id int64, err error) {
	pwd := utils.EncryptPassword(password)
	sqlCmd := "insert into user(nick_name, mobile, password) values (?, ?, ?)"
	ret, err := db.Exec(sqlCmd, nickname, mobile, pwd)
	if err != nil {
		return 0, err
	}
	return ret.LastInsertId()
}

//ModifyUserInfo 修改用户信息
func ModifyUserInfo(id int64, nickname, gender, birthday string) error {
	sqlCmd := `update user set nick_name=?, gender=?, birthday=? where id=?`
	_, err := db.Exec(sqlCmd, nickname, gender, birthday, id)
	return err
}

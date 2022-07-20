package main

import (
	"log"
	"luke544187758/goods-srv/dao/mysql"
	"luke544187758/goods-srv/settings"
)

func main() {

	cfg := new(settings.MySQLConfig)
	cfg.Host = "127.0.0.1"
	cfg.Port = 13306
	cfg.User = "root"
	cfg.Password = "544187758"
	cfg.DB = "mall-goods-srv"
	cfg.MaxIdleConns = 50
	cfg.MaxOpenConns = 200
	if err := mysql.Init(cfg); err != nil {
		log.Fatal(err)
	}

}

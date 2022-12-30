package main

import (
	"fmt"
	"luke544187758/order-srv/dao/mysql"
	"luke544187758/order-srv/models"
	"luke544187758/order-srv/pkg/snowflake"
	"luke544187758/order-srv/settings"
	"time"
)

func main() {
	cfg := &settings.MySQLConfig{
		Port:         13306,
		MaxOpenConns: 200,
		MaxIdleConns: 50,
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "544187758",
		DB:           "mall-order-srv",
	}
	if err := mysql.Init(cfg); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	fmt.Println("mysql init success...")

	if err := snowflake.Init("2022-03-20", 1); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	items := make([]*models.ShoppingCart, 0)
	items = append(items, &models.ShoppingCart{
		ID:         snowflake.GenID(),
		User:       1,
		Goods:      316769164070912,
		Nums:       1,
		Checked:    true,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	time.Sleep(time.Microsecond * 100)

	items = append(items, &models.ShoppingCart{
		ID:         snowflake.GenID(),
		User:       1,
		Goods:      316972973690880,
		Nums:       2,
		Checked:    true,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	time.Sleep(time.Microsecond * 100)

	items = append(items, &models.ShoppingCart{
		ID:         snowflake.GenID(),
		User:       1,
		Goods:      318763161358336,
		Nums:       4,
		Checked:    true,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	time.Sleep(time.Microsecond * 100)

	items = append(items, &models.ShoppingCart{
		ID:         snowflake.GenID(),
		User:       1,
		Goods:      318786989199360,
		Nums:       3,
		Checked:    true,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	time.Sleep(time.Microsecond * 100)

	items = append(items, &models.ShoppingCart{
		ID:         snowflake.GenID(),
		User:       1,
		Goods:      318893650350080,
		Nums:       2,
		Checked:    true,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	for _, v := range items {
		if err := mysql.CreateCartItem(v); err != nil {
			fmt.Printf("[%d] 创建失败, %v\n", v.Goods, err)
		}
	}

}

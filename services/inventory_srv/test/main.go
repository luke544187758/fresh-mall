package main

import (
	"fmt"
	"luke544187758/inventory-srv/dao/mysql"
	"luke544187758/inventory-srv/models"
	"luke544187758/inventory-srv/pkg/snowflake"
	"luke544187758/inventory-srv/settings"
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
		DB:           "mall-inventory-srv",
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

	invs := make([]*models.Inventory, 0)
	invs = append(invs, &models.Inventory{
		ID:         snowflake.GenID(),
		Goods:      316769164070912,
		Stocks:     100,
		Version:    0,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	time.Sleep(time.Microsecond * 100)
	invs = append(invs, &models.Inventory{
		ID:         snowflake.GenID(),
		Goods:      316972973690880,
		Stocks:     100,
		Version:    0,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	time.Sleep(time.Microsecond * 100)
	invs = append(invs, &models.Inventory{
		ID:         snowflake.GenID(),
		Goods:      318763161358336,
		Stocks:     100,
		Version:    0,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	time.Sleep(time.Microsecond * 100)
	invs = append(invs, &models.Inventory{
		ID:         snowflake.GenID(),
		Goods:      318786989199360,
		Stocks:     100,
		Version:    0,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	time.Sleep(time.Microsecond * 100)
	invs = append(invs, &models.Inventory{
		ID:         snowflake.GenID(),
		Goods:      318893650350080,
		Stocks:     100,
		Version:    0,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	time.Sleep(time.Microsecond * 100)
	invs = append(invs, &models.Inventory{
		ID:         snowflake.GenID(),
		Goods:      318946599243776,
		Stocks:     100,
		Version:    0,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	time.Sleep(time.Microsecond * 100)
	invs = append(invs, &models.Inventory{
		ID:         snowflake.GenID(),
		Goods:      319094700118016,
		Stocks:     100,
		Version:    0,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	time.Sleep(time.Microsecond * 100)
	invs = append(invs, &models.Inventory{
		ID:         snowflake.GenID(),
		Goods:      326308630368256,
		Stocks:     100,
		Version:    0,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	time.Sleep(time.Microsecond * 100)
	invs = append(invs, &models.Inventory{
		ID:         snowflake.GenID(),
		Goods:      338987021504512,
		Stocks:     100,
		Version:    0,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	time.Sleep(time.Microsecond * 100)
	invs = append(invs, &models.Inventory{
		ID:         snowflake.GenID(),
		Goods:      339213220319232,
		Stocks:     100,
		Version:    0,
		IsDeleted:  false,
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	})
	for _, v := range invs {
		if err := mysql.CreateInventory(v); err != nil {
			fmt.Printf("[%d] 创建失败, %v\n", v.Goods, err)
		}
	}
}

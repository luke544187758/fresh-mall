package redis

import (
	"context"
	"errors"
	"time"
)

// 重试次数
var retryTimes = 100

// 重试频率
var retryInterval = time.Millisecond * 100

// 锁的默认过期时间
var expiration = time.Second * 15

//LockerAcquire 加锁
func LockerAcquire(ctx context.Context, key string, value interface{}) error {
	//尝试加锁
	ok_lock, err := rdb.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		return err
	}
	// 加锁失败，重试
	if !ok_lock {
		ok_retry, err := retry(ctx, key, value)
		if err != nil {
			return err
		}
		if !ok_retry {
			return errors.New("server unavailable, try again later")
		}
	}
	// 加锁成功,开启守护线程
	go watchDog(ctx, key)
	return nil
}

//LockerRelease 释放锁
func LockerRelease(ctx context.Context, key string, value interface{}) error {
	lua := `
-- 如果当前值与锁值一致,删除key
if redis.call('GET', KEYS[1]) == ARGV[1] then
	return redis.call('DEL', KEYS[1])
else
	return 0
end
`
	scriptKeys := []string{key}

	_, err := rdb.Eval(ctx, lua, scriptKeys, value).Result()
	if err != nil {
		return err
	}

	return nil
}

//watchDog 守护线程
func watchDog(ctx context.Context, key string) {
	for {
		select {
		// 业务完成
		case <-ctx.Done():
			return
			// 业务未完成
		default:
			// 自动续期
			rdb.PExpire(ctx, key, expiration)
			// 继续等待
			time.Sleep(expiration * 2 / 3)
		}
	}
}

//retry 重试
func retry(ctx context.Context, key string, value interface{}) (bool, error) {
	i := 1
	for i <= retryTimes {
		set, err := rdb.SetNX(ctx, key, value, expiration).Result()

		if err != nil {
			return false, err
		}

		if set == true {
			return true, nil
		}

		time.Sleep(retryInterval)
		i++
	}
	return false, nil
}

package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"time"
)

// 通过redis获取分布式锁
// @Param lockName 锁的名字
// @Param acquireTimeout  等待的时长
// @Param lockTimeout 锁的有效期
func RedisGetLock(lockName string, acquireTimeout time.Duration, lockTimeout time.Duration) (bool, error) {
	code := time.Now().String()
	endTime := time.Now().Add(acquireTimeout).UnixNano()
	for time.Now().UnixNano() <= endTime {
		success, err := RedisClient.SetNX(lockName, code, lockTimeout).Result()
		if err != nil {
			return false, err
		} else if success {
			return true, nil
		} else if RedisClient.TTL(lockName).Val() == -1 {
			RedisClient.Expire(lockName, lockTimeout)
		} else if !success {
			//锁已被占用
			return false, nil
		}
		time.Sleep(time.Millisecond)
	}
	return false, fmt.Errorf("timeout")
}

func RedisReleaseLock(lockName string) bool {
	txf := func(tx *redis.Tx) error {
		if _, err := tx.Get(lockName).Result(); err != nil && err != redis.Nil {
			return err
		} else {
			_, err := tx.Pipelined(func(pipe redis.Pipeliner) error {
				pipe.Del(lockName)
				return nil
			})
			return err
		}
	}

	for {
		if err := RedisClient.Watch(txf, lockName); err == nil {
			return true
		} else if err == redis.TxFailedErr {
			logrus.Errorf("ERROR: watch key is modified, try to release lock. err : ", err)
		} else {
			logrus.Errorf("ERROR: ", err)
			return false
		}
	}
}

func IsExpire(key string) bool {
	// 过期或不存在都返回这个
	return RedisClient.TTL(key).Val() == -2 * time.Second || RedisClient.TTL(key).Val() == -1 * time.Second
}

func LockKey(key string, lockLeaseTime time.Duration) error {
	_, err := RedisClient.SetNX(key, "1", lockLeaseTime).Result()
	if err != nil {
		return err
	}
	logrus.Info("INFO: lock key %s success", key)
	return nil
}

func RenewExpiration(key string, lockLeaseTime time.Duration) error {
	success, err := RedisClient.Expire(key, lockLeaseTime).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	// 不存在或者没有续期成功
	if !success && IsExpire(key) {
		log.Info("INFO: instance %s retry get lock", key)
		err := LockKey(key, lockLeaseTime)
		if err != nil {
			return err
		}
		success = true
	}
	return nil
}

func DelKey(key string) error {
	_, err := RedisClient.Del(key).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}
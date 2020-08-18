package instance

import (
	"bryson.foundation/kbuildresource/async"
	"bryson.foundation/kbuildresource/cache"
	"bryson.foundation/kbuildresource/common"
	"bryson.foundation/kbuildresource/utils"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	lockLeasTime = 6 * time.Second // 每次续期6s
	renewInterval = 3 * time.Second // 每3秒续期一次
)

var (
	instanceKeyOfSelf = ""
	distributeLockKey = cache.GenMetaDistributeKey(common.BuildJobPrefix)
	instanceNameListKey = cache.GenInstanceNameListKey(common.BuildJobPrefix)
)

type LivenessProbe interface {
	prepare() // 准备工作
	isLive() bool
	clear() // 清理工作
}

type Instance interface {
	LivenessProbe
	Name() string // 返回名字
	StartUp() //启动应用实例
	Shutdown() //关闭应用
	postStart() //应用实例启动之后的操作
	preStop() // 应用实例关闭之前的操作
}

type instanceWithRedis struct {
	name string
	stopCh chan struct{} // 程序自动调用关闭实例
	signalCh chan os.Signal // 接收到syscall.SIGINT和syscall.SIGTERM 信号量关闭
	liveCh chan struct{} // 程序存活
	requestController *async.RequestController
}

var (
	BeeInstance Instance
)

func init() {
	BeeInstance = newInstance()
}

func newInstance() Instance {
	instanceName := utils.CreateRandomString(8)
	instance := &instanceWithRedis{
		name:              instanceName,
		stopCh:            make(chan struct{}),
		signalCh:          make(chan os.Signal),
		liveCh:            make(chan struct{}),
		requestController: async.NewRequestController(),
	}
	instanceKeyOfSelf = cache.GenInstanceKey(common.BuildJobPrefix, instance.name)
	return instance
}

func (instance *instanceWithRedis) postStart() {
	signal.Notify(instance.signalCh, syscall.SIGINT, syscall.SIGTERM)

	// 启动requestController
	go instance.requestController.StartUp()

	// 接收其他崩溃的instance的job任务
	// 获取同步锁
	go func() {
		logrus.Info("INFO: start takeover job of other instance")
		isNeed := true // 默认需要接收
		// 尝试获取同步锁
		for i := 0; i < 3; i++ {
			isSuccess, err := cache.RedisGetLock(distributeLockKey, 10 *time.Second, 2 * time.Minute)
			if err != nil {
				logrus.Error("ERROR: lock distribute key failed, err: ", err)
				continue
			}
			// 锁已被占用
			if !isSuccess {
				logrus.Info("INFO: lock distribute key faield, which is already be lock")
				isNeed = false
			}
			break
		}
		if !isNeed {
			return
		}
		//获取锁成功，开始检索任务，并把自己添加到实例列表里面
		instanceNameListVal, err := cache.RedisClient.Get(instanceNameListKey).Result()
		instanceNameList := make([]string, 0)
		if err != nil && err != redis.Nil {
			logrus.Errorf("ERROR: get instanceNameListKey failed, error: ", err)
			return
		} else if err == redis.Nil{
			//如果不存在，不操作
		} else {
			err = json.Unmarshal([]byte(instanceNameListVal), &instanceNameList)
			if err != nil {
				logrus.Error("ERROR: unmarshal failed")
				return
			}
		}
		liveInstances := make([]string, 0)
		for _, instanceName := range instanceNameList {
			if !checkLive(instanceName) {
				err := instance.requestController.TakeOverRequest(instanceName, instance.name)
				if err != nil {
					logrus.Error("ERROR: ", err)
				}
			} else {
				liveInstances = append(liveInstances, instanceName)
			}
		}
		// 把自己添加到instance-list里
		liveInstances = append(liveInstances, instance.name)
		instanceNameListJsonData, err := json.Marshal(liveInstances)
		if err != nil {
			logrus.Error("ERROR: ", err)
			return
		}
		err = cache.RedisClient.Set(instanceNameListKey, instanceNameListJsonData,0).Err()
		if err != nil {
			logrus.Error("ERROR: update instanceNameList failed")
			return
		}

		//释放锁
		isRelease := cache.RedisReleaseLock(distributeLockKey)
		if isRelease {
			logrus.Info("INFO: release distribute lock success")
		} else {
			logrus.Info("INFO: release distribute lock failed")
		}
	}()
}

// 做服务优雅停机
func (instance *instanceWithRedis) preStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()
	if err := beego.BeeApp.Server.Shutdown(ctx); err != nil {
		logrus.Fatal("ERROR: server force to shutdown: ", err)
	}
	instance.requestController.Shutdown()
}

func (instance *instanceWithRedis) Name() string {
	return instance.name
}

// 提供外部调用，用于关闭应用
func (instance *instanceWithRedis) Shutdown() {
	close(instance.stopCh)
}

func (instance *instanceWithRedis) StartUp() {
	// 做存活性探针的相关工作
	instance.prepare()
	defer instance.clear()

	instance.postStart()
	// 监控各种信号
	for {
		select {
		// 监控stopCh信号， 执行程序退出 前处理操作，是结束实例的唯一入口
		case <-instance.stopCh:
			logrus.Info("INFO: shutting down the server because instance.Shutdown()")
			instance.preStop()
			logrus.Info("INFO: finish preStop")
			return
		case signalVal := <- instance.signalCh:
			logrus.Infof("INFO: shutting down the sever beacause signalCh %v,", signalVal)
			instance.Shutdown()
		case <-instance.liveCh:
			logrus.Infof("INFO: %s instance is live ", instance.name)
		}
	}
}

// 基于redis 来实现存活性验证
func (instance *instanceWithRedis) prepare() {
	// 不断续期，直到死亡
	err := cache.LockKey(instanceKeyOfSelf, lockLeasTime)
	if err != nil {
		logrus.Errorf("ERROR: instance %s get lock failed, err: ", instance.name, err)
		return
	}
	ticker := time.NewTicker(renewInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case t := <-ticker.C:
				err := cache.RenewExpiration(instanceKeyOfSelf, lockLeasTime)
				if err != nil {
					logrus.Infof("INFO: instance %s renew lock failed at %s", instance.name, t)
					continue
				}
				instance.liveCh <- struct{}{}
			case <-instance.stopCh:
				return

			}
		}
	}()
}

func (instance *instanceWithRedis) clear() {
	logrus.Info("INFO: delete redis key ", instanceKeyOfSelf)
	err := cache.DelKey(instanceKeyOfSelf)
	if err != nil {
		logrus.Errorf("INFO: clear instanceKeyOfSelf %s failed, err :", instanceKeyOfSelf, err)
		return
	}
	logrus.Infof("INFO: clear instanceKeyOfSelf %s success", instanceKeyOfSelf)
}

func (instance *instanceWithRedis) isLive() bool {
	return !cache.IsExpire(instanceKeyOfSelf)
}

func checkLive(instanceName string) bool {
	key := cache.GenInstanceKey(common.BuildJobPrefix, instanceName)
	return !cache.IsExpire(key)
}
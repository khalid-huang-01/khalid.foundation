package instance

import (
	"bryson.foundation/kbuildresource/async"
	"bryson.foundation/kbuildresource/cache"
	"bryson.foundation/kbuildresource/common"
	"bryson.foundation/kbuildresource/utils"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/sirupsen/logrus"
	"math/rand"
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
	// 确保把自己添加到实例列表中
	result := RetrieveAccessOfUpdateInstanceNameList()
	if !result {
		logrus.Infof("INFO: RetrieveAccessOfUpdateInstanceNameList failed")
		return
	}
	instanceNameList, err := cache.GetInstanceNameList()
	if err != nil {
		return
	}
	instanceNameList = append(instanceNameList, instance.name)
	instanceNameListJsonData, err := json.Marshal(instanceNameList)
	if err != nil {
		logrus.Error("ERROR: ", err)
		return
	}
	err = cache.RedisClient.Set(instanceNameListKey, instanceNameListJsonData,0).Err()
	if err != nil {
		logrus.Error("ERROR: update instanceNameList failed")
		return
	}
	ReturnAccessOfUpdateInstanceNameList()

	// 启动requestController
	go instance.requestController.StartUp()

	// 接收其他崩溃的instance的job任务
	// 获取同步锁
	// 如果有多个同时在获取，只有一个成功即可
	go func() {
		logrus.Info("INFO: start takeover job of other instance")
		isSuccess := RetrieveAccessOfTakeOver()
		if !isSuccess {
			return
		}
		defer ReturnAccessOfTakeOver()


		//获取锁成功，开始检索任务，
		instanceNameList, err := cache.GetInstanceNameList()
		if err != nil {
			return
		}

		liveInstances := make([]string, 0)
		for _, instanceName := range instanceNameList {
			if !checkLive(instanceName) {
				go func(instanceName string) {
					err := instance.requestController.TakeOverRequest(instanceName, instance.name)
					if err != nil {
						logrus.Error("ERROR: ", err)
					}
				}(instanceName)
			} else {
				liveInstances = append(liveInstances, instanceName)
			}
		}
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

// 这个函数需要传入一个finishCh来通知外部调用者，内部已经初始化完成
func (instance *instanceWithRedis) StartUp() {
	//监听信号
	signal.Notify(instance.signalCh, syscall.SIGINT, syscall.SIGTERM)

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
	_, err := cache.LockKey(instanceKeyOfSelf, lockLeasTime)
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

// 获取分布式锁，如果占用需要等待
func RetrieveAccessOfUpdateInstanceNameList() bool {
	result := true
	var err error
	for i := 0; i < 10; i++ {
		result, err = cache.LockKey(distributeLockKey, time.Minute)
		// 网络问题进行重试
		if err != nil {
			logrus.Error("ERROR: ", err)
			continue
		} else if !result {
			// 获取不到，重新随机等待，再重新获取
			time.Sleep(time.Duration(200 + rand.Int63n(1000)) * time.Millisecond)
			continue
		}
		// 成功直接break
		break
	}
	return result
}

func ReturnAccessOfUpdateInstanceNameList() bool {
	return releaseDistributeKey()
}

// 获取 分布式锁，非网络原因，直接返回结果
func RetrieveAccessOfTakeOver() bool {
	result := true
	var err error
	for i := 0; i < 3; i++ {
		result, err = cache.LockKey(distributeLockKey, time.Minute)
		// 网络问题进行重试
		if err != nil {
			continue
		}
		break
	}
	return result
}

func ReturnAccessOfTakeOver() bool {
	return releaseDistributeKey()
}

func releaseDistributeKey() bool {
	return cache.RedisReleaseLock(distributeLockKey)
}
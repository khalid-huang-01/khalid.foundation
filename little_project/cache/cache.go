// 参考java的guava 本地缓存管理器
//适用于
//	愿意消耗一些内存空间来提升速度
//	预料至某些键会被查询一次以上
//	缓存中存放的数据问题不会超出内存容
//功能有
//	值有过期时间，默认会自动刷新，刷新函数由外部程序自定义
//	带有容量，避免占用过多内存
//	超过阀值会自动清理
//	清理可以按LRU或者FIFO策略
//	LRU通过updateTime来，最近使用的不淘汰

package cache

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

type Item struct {
	key      interface{}
	value    interface{}
	addTime  time.Time
	readTime time.Time
}

var clearThreshold float32 = 0.85 // 定义清理阈值为size占用超过容量85%

type ItemSlice []*Item

// 按照CacheItem.readTime 从前到后排序
func (c ItemSlice) Len() int {
	return len(c)
}

func (c ItemSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ItemSlice) Less(i, j int) bool {
	return c[j].readTime.After(c[i].readTime)
}

func (c Item) isExpire(expire time.Duration) bool {
	if expire.Microseconds() == 0 {
		return false
	}
	return c.addTime.Add(expire).Before(time.Now())
}

const defaultGetTimeout = time.Duration(30 * time.Second) // 网关超时时间，用于使用loader函数刷新键值的时候

type Manager struct {
	sync.RWMutex

	capacity int // cache中键值对总数上限，超过总数，按LRU驱逐旧键值
	expire time.Duration //键值的过期时间 // 对于过期的键值当前没有怎么处理
	renew bool

	writeCount int64 // 增加键值的计数器，永远只增长
	cleanCount int64 // 运行清除任务次数
	TotalCleanTime time.Duration // 运行清除任务的总时间

	cacheMap map[interface{}]*Item
	loader func(interface{}) (interface{}, error) // 刷新键值的函数

	cleanTaskChan chan bool // 键值清理通知标识
}

//创建一个LRU缓存管理器，key为string 主要是为了兼容普遍场景
func NewCacheManager(capacity int, expire time.Duration, loader func(string) (interface{}, error)) (*Manager, error) {
	return newCacheManager(capacity, expire, true, func(key interface{}) (interface{}, error) {
		return loader(key.(string))
	})
}

//创建一个LRU缓存管理器, key 为任意可比较的interface{}
//参数：capacity: >0 为键值容量，为0表示不控制键值大小
//		expire: 键值过期时间，为0表示不做过期处理，带过期处理的时候按照LRU算法进行
//      loader: 加载键值的函数
func NewCacheManagerExtKey(capacity int, expire time.Duration, loader func(interface{}) (interface{}, error)) (*Manager, error) {
	return newCacheManager(capacity, expire, true, loader)
}
// 创建一个LRU缓存管理器，key为string
func NewCacheManagerFIFO(capacity int, expire time.Duration, loader func(string) (interface{}, error)) (*Manager, error) {
	return newCacheManager(capacity, expire, false, func(key interface{}) (interface{}, error) {
		return loader(key.(string))
	})
}

//创建一个LRU缓存管理器
//参数：capacity: >0 为键值容量，为0表示不控制键值大小
//		expire: 键值过期时间，为0表示不做过期处理
//      loader: 加载键值的函数
func NewCacheManagerFIFOExtKey(capacity int, expire time.Duration, loader func(interface{}) (interface{}, error)) (*Manager, error) {
	return newCacheManager(capacity, expire, false, loader)
}

func newCacheManager(capacity int , expire time.Duration, renew bool, loader func(interface{}) (interface{}, error)) (*Manager, error) {
	log.Printf("INFO: create cache manager capacity %d, expire %v", capacity, expire)

	if capacity < 0 {
		return nil, fmt.Errorf("invalid capacity %d", capacity)
	}

	m := Manager{}
	m.loader = loader
	m.capacity = capacity
	m.expire = expire
	m.renew = renew
	m.cacheMap = make(map[interface{}]*Item, capacity/2)
	m.cleanTaskChan = make(chan bool, 1000)
	go m.clean()
	return &m, nil
}



// 获取key指定的键值，如果在缓存中直接命中，如果不在缓存中，则调用loader函数获取
// timeout 指定调用load函数时个的超时时间，如果为0，则使用缺省超时时间30s
func (c *Manager) Get(key string, timeout time.Duration) (interface{}, error) {
	return c.GetExtKey(key, timeout)
}

func (c *Manager) GetExtKey(key interface{}, timeout time.Duration) (interface{}, error) {
	c.RLock() // 读锁
	oldCacheValue, okay := c.cacheMap[key]
	if okay {
		//存在且未过期，则直接返回
		if !oldCacheValue.isExpire(c.expire) {
			if c.renew {
				oldCacheValue.readTime = time.Now() // 设置最新使用时间，用于LRU淘汰
			}
			c.RUnlock()
			return oldCacheValue.value, nil
		}
	}
	c.RUnlock()

	// 缓存未命中，或者已过期
	var value interface{}
	readChan := make(chan bool, 1) // chan里值 代表获取的成功或者失败状态
	if timeout <= 0 {
		timeout = defaultGetTimeout
	}
	go func() {
		defer close(readChan)
		//尝试loader
		start := time.Now()
		interval := 10 * time.Millisecond // 起始时间为10ms
		var err error
		for {
			value, err = c.loader(key)
			if err == nil {
				readChan <- true
				return
			}
			if interval > time.Second {
				interval = time.Second
			}
			if time.Now().After(start.Add(timeout + interval)) {
				//已经重试并彻底超时了
				readChan <- false
				return
			}
			time.Sleep(interval)
			interval *= 2
		}
	}()
	//等待新值到来或者超时
	select {
	case success, okay := <-readChan:
		if success && okay {
			//设置并返回新值
			c.set(key, value)
			return value, nil
		}
	case <-time.After(timeout):
	}

	if oldCacheValue != nil {
		return oldCacheValue.value, nil
	}

	return nil, fmt.Errorf("cannot find key %s", key)
}

func (c *Manager) set(key interface{}, value interface{}) {
	c.Lock()
	defer c.Unlock()

	if c.Size() >= c.capacity {
		//已经无法写入，强制运行
		c.doClean(true)
	} else if c.needClean() {
		c.cleanTaskChan <- false
	}

	c.cacheMap[key] = &Item{
		key:      key,
		value:    value,
		addTime:  time.Now(),
		readTime: time.Now(),
	}
	c.writeCount++
}

func (c *Manager) Size() int {
	return len(c.cacheMap)
}

func (c *Manager) needClean() bool {
	return c.Size() >= int(float32(c.capacity) * clearThreshold)
}

// 清除到size 小于capacity的85%
// inWriteLock用于判断是否调用方已经加了锁，因为这个清理是需要加锁处理的
func (c *Manager) doClean(inWriteLock bool) {
	if c.Size() == 0 {
		return
	}
	if !inWriteLock {
		c.Lock()
		defer c.Unlock()
	}
	startTime := time.Now()
	c.cleanCount++

	var cacheSlice ItemSlice = make([]*Item, 0, c.Size())
	for _, value := range c.cacheMap {
		cacheSlice = append(cacheSlice, value)
	}
	sort.Sort(cacheSlice)
	cleanCount := int(float32(c.capacity) * clearThreshold) - c.Size()
	if cleanCount <= 0 {
		// 起码清理一个
		cleanCount = 1
	}
	if cleanCount > len(cacheSlice) {
		cleanCount = len(cacheSlice)
	}
	for i := 0; i < cleanCount; i++ {
		delete(c.cacheMap, cacheSlice[i].key)
	}
	c.TotalCleanTime += time.Now().Sub(startTime)
}

func (c *Manager) Close() {
	log.Print("INFO: cache manager will close")
	close(c.cleanTaskChan)
}

// 清理过期键值
func (c *Manager) clean() {
	for {
		var period bool
		overCapacity, okay := <-c.cleanTaskChan
		if !okay {
			log.Print("INFO: Receive cache manager close notification")
			return
		}
		period = !overCapacity
		// 非阻塞的读出所有缓存的清理任务项，只执行其中一个任务，优先执行非阻塞的[这里有问题的，只能读取一次, 实际应该是全读取出来，然后只执行一次]
		select {
		case overCapacity, okay = <-c.cleanTaskChan:
			if !okay {
				log.Print("INFO: Receive cache manager close notification")
				return
			}
			if !period {
				period = !overCapacity
			}
		default:
		}
		c.doClean(false)
	}
}

func (c *Manager) keys() []interface{} {
	c.RLock()
	defer c.RUnlock()

	var keys []interface{}
	for key, _ := range c.cacheMap {
		keys = append(keys, key)
	}
	return keys
}

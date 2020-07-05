package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCacheManagerFIFO(t *testing.T) {
	var capacity int = 100
	var expire time.Duration = 10 * time.Second
	manager, err := NewCacheManagerFIFO(capacity,expire, loader)
	if err != nil {
		fmt.Println("Err: ", err)
		return
	}
	// 测试
	for i := 0; i < 100; i++{
		key := "fan"
		value, err := manager.Get(key, 3 * time.Second)
		if err != nil {
			fmt.Println("ERROR: ", err)
			continue
		}
		fmt.Printf("Key: %s, Value: %s \n", key, value)
		time.Sleep(time.Second)
	}
}

func loader(key string) (interface{}, error) {
	fmt.Println("重新获取")
	time.Sleep(time.Second)
	return key+":value", nil
}
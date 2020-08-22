package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"testing"
	"time"
)

func TestLock(t *testing.T){
	c, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	l := zk.NewLock(c, "/lock", zk.WorldACL(zk.PermAll))
	err = l.Lock()
	time.AfterFunc(4 * time.Second, func() {
		l := zk.NewLock(c, "/lock", zk.WorldACL(zk.PermAll))
		err = l.Lock()
		t.Log("get lock")
		err = l.Unlock()
	})
	if err != nil {
		t.Fatal("Error: ", err)
	}
	t.Log("lock succ, do your business logic")

	time.Sleep(time.Second * 10)

	// do some thing
	l.Unlock()
	t.Log("unlock succ, finish business logic")
	time.Sleep(time.Second * 5)
}
package design

import (
	"container/list"
	"fmt"
	"testing"
)

// 146
// LRU数据需要移动，移动需要O(1)达到，需要使用双向链表
// 查找需要O(1),需要哈希表
// => 数据存储层使用双向链表，索引层使用哈希表
type LRUCache struct {
	capacity int
	index map[int]*list.Element
	store *list.List
}

type Pair struct {
	key int
	value int
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		index:    make(map[int]*list.Element),
		store:    list.New(),
	}
}


func (this *LRUCache) Get(key int) int {
	node, ok := this.index[key]
	if !ok {
		return -1
	}
	// 移动到前面
	this.store.MoveToFront(node)

	return node.Value.(Pair).value
}


func (this *LRUCache) Put(key int, value int)  {
	node, ok := this.index[key]
	if ok {
		node.Value = Pair{
			key:   key,
			value: value,
		}
		this.store.MoveToFront(node)
	} else {
		// 如果达标要淘汰
		if this.store.Len() >= this.capacity {
			backNode := this.store.Back()
			this.store.Remove(backNode)
			delete(this.index, backNode.Value.(Pair).key)
		}
		// 加入新的kv
		newNode := this.store.PushFront(Pair{
			key:   key,
			value: value,
		})
		this.index[key] = newNode
	}
}


/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */

func TestList(t *testing.T)  {
	doubleList := list.New()
	doubleList.PushFront(1)
	doubleList.PushFront(2)

	cur := doubleList.Front()
	for cur != nil {
		fmt.Println(cur.Value)
		cur = cur.Next()
	}
}

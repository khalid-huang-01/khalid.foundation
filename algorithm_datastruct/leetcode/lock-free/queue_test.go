package lock_free

import (
	"fmt"
	"sync/atomic"
	"testing"
	"unsafe"
)

type Node struct {
	value interface{}
	next unsafe.Pointer
}

type LKQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

// 虚拟头部节点
func NewLKQueue() *LKQueue {
	n := unsafe.Pointer(&Node{})
	return &LKQueue{head: n, tail: n}
}

func (q *LKQueue) Enqueue(v interface{}) {
	n := &Node{value: v}
	// cas不断测试
	for {
		tail := load(&q.tail)
		next := load(&tail.next)
		if tail == load(&q.tail) {
			if next == nil {
				if cas(&tail.next, next, n) {
					cas(&q.tail, tail, n)
					return
				}
			} else {
				// 如果next已经不是null了，表明有新的插入了，tail就往下走
				// 往下走的条件是自己这里的就是最新的tail
				cas(&q.tail, tail, next)
			}
		}
	}
}

func (q *LKQueue) Dequeue() interface{} {
	for {
		head := load(&q.head)
		tail := load(&q.tail)
		next := load(&head.next)
		if head == load(&q.head) {
			if head == tail {
				if next == nil {
					return nil
				}
				cas(&q.tail, tail, next)
			} else {
				v := next.value
				if cas(&q.head, head, next) {
					return v
				}
			}
		}
	}
}

func load(p *unsafe.Pointer) (n *Node) {
	return (*Node)(atomic.LoadPointer(p))
}

func cas(p *unsafe.Pointer, old, new *Node) bool {
	return atomic.CompareAndSwapPointer(
		p, unsafe.Pointer(old), unsafe.Pointer(new))
}

func TestLFQueue(t *testing.T)  {
	//入队
	q := NewLKQueue()
	q.Enqueue(3)
	q.Enqueue(1)
	q.Enqueue(2)
	// 出队
	v := q.Dequeue()
	fmt.Println(v.(int))

	v = q.Dequeue()
	fmt.Println(v.(int))
	v = q.Dequeue()
	fmt.Println(v.(int))
}
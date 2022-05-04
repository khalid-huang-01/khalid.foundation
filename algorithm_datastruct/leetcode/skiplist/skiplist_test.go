package skiplist

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const (
	MaxLevel = 32
)

type Node struct {
	val   int
	level int
	next  []*Node
}

func NewNode(val, level int) *Node {
	return &Node{
		val:   val,
		level: level,
		next:  make([]*Node, level),
	}
}

type Skiplist struct {
	level int
	head  *Node
}

func Constructor() Skiplist {
	rand.Seed(time.Now().UnixNano())
	return Skiplist{
		level: MaxLevel,
		head:  NewNode(0, MaxLevel), //虚拟头节点
	}
}

func (s *Skiplist) Search(target int) bool {
	prev := s.head
	// 从最高层往下找
	for i := s.level - 1; i >= 0; i-- {
		// 同层遍历
		cur := prev.next[i]
		for cur != nil {
			if cur.val == target {
				return true
			} else if cur.val < target {
				// 本层继续
				prev = cur
				cur = prev.next[i]
			} else {
				// 不在本层了, prev保持不动，prev是下一层遍历的起点
				break
			}
		}
	}
	return false
}

func (s *Skiplist) Add(num int) {
	// fmt.Println("add ", num)
	prevs := make([]*Node, s.level)
	prev := s.head
	for i := s.level - 1; i >= 0; i-- {
		prevs[i] = prev // 先添加一个默认值，后续不断更新，兼容虚拟头节点（主要是为了前驱是虚拟头节点的情况）
		cur := prev.next[i]
		for cur != nil {
			// 等于与大于或者是最后一个就可以保存前驱返回了
			if cur.val >= num {
				prevs[i] = prev
				break // 跳出本层遍历
			}
			// 如果已经是最后一个了，还没有找到比这个大的，那么最后一个元素就是前驱
			if cur.next[i] == nil {
				prevs[i] = cur
				prev = cur // 下一轮遍历的起点
				break
			}
			// 继续遍历本层
			prev = cur
			cur = prev.next[i]
		}
	}

	// 插入节点
	level := s.randomLevel()
	// fmt.Println("level: ", level)
	// fmt.Printf("%+v", prevs)
	node := NewNode(num, level)
	for i := 0; i < level; i++ {
		node.next[i] = prevs[i].next[i]
		prevs[i].next[i] = node
	}

}

func (s *Skiplist) Erase(num int) bool {
	// fmt.Println("erase ", num)
	// 找到位置
	var targetNode *Node
	prevs := make([]*Node, s.level)
	prev := s.head
	for i := s.level - 1; i >= 0; i-- {
		prevs[i] = prev // 先添加一个默认值，后续不断更新，兼容虚拟头节点
		cur := prev.next[i]
		for cur != nil {
			if cur.val == num {
				targetNode = cur
				prevs[i] = prev // 保存前驱
				break           // 进入下一层
			} else if cur.val < num {
				prev = cur
				cur = prev.next[i]
			} else {
				break // 进入下一层
			}
		}
	}
	// 修改指向
	if targetNode == nil {
		return false
	}
	for i := 0; i < targetNode.level; i++ {
		prevs[i].next[i] = targetNode.next[i]
	}
	return true
}

func (s *Skiplist) print(level int) {
	cur := s.head
	for cur != nil {
		fmt.Printf("%d->", cur.val)
		cur = cur.next[level]
	}
	fmt.Println("")
}

func (s *Skiplist) randomLevel() int {
	// return 1
	i := 1
	for ; i < s.level; i++ {
		if rand.Intn(2) == 0 { // [0,2)
			return i
		}
	}
	return i
}

/**
 * Your Skiplist object will be instantiated and called as such:
 * obj := Constructor();
 * param_1 := obj.Search(target);
 * obj.Add(num);
 * param_3 := obj.Erase(num);
 */

func TestXxx(t *testing.T) {
	s := Constructor()
	s.Add(0)
	s.Add(5)
	s.Add(2) // 插在中间的时候会超时
	s.Add(2) // 插在中间的时候会超时
	s.Add(2) // 插在中间的时候会超时
	s.Add(1)
	t.Log(s.Search(2))
	// t.Log(s.Search(2))
	s.print(0)
	// s.print(0)
	// t.Log(s.Search(0))
	// // s.Erase(5)
	// t.Log(s.Search(2))
	// s.print(0)
	// t.Log(s.Search(3))
	// t.Log(s.Search(2))
	// s.Erase(2)
	// t.Log(s.Search(2))
}

package recursion

import (
	"fmt"
	"testing"
)

// https://juejin.cn/post/6844903910386171912
type ListNode struct {
	Val int
	Next *ListNode
}
func reverseKGroup(head *ListNode, k int) *ListNode {
	// 递归终止条件：判断当前的长度够不够凑齐k
	tmp := head
	tmpK := 0
	for tmp != nil && tmpK < k {
		tmp = tmp.Next
		tmpK += 1
	}
	if tmpK < k {
		return head
	}

	// 做当前的运算：对当前的进行翻转
	cur := head
	var pre, next *ListNode
	// 总共需要翻转k次
	for i := 1; i <= k; i++ {
		next = cur.Next // 保留后面的元素
		cur.Next = pre // 修改新指向
		pre = cur // 指标推进
		cur = next // 指标推进
	}
	newHead := pre
	newTail := head // 一开始的头就是现在的尾巴
	// 进行递归：把当前的和后面的连接起来（递归）
	newTail.Next = reverseKGroup(next, k)

	// 返回结果
	return newHead
}

func TestReverseKGroup(t *testing.T) {
	v8 := &ListNode{
		Val: 8,
	}
	v7 := &ListNode{
		Val: 7,
		Next: v8,
	}
	v6 := &ListNode{
		Val:  6,
		Next: v7,
	}
	v5 := &ListNode{
		Val: 5,
		Next: v6,
	}
	v4 := &ListNode{
		Val: 4,
		Next: v5,
	}
	v3 := &ListNode{
		Val:  3,
		Next: v4,
	}
	v2 := &ListNode{
		Val: 2,
		Next: v3,
	}
	head := &ListNode{
		Val: 1,
		Next: v2,
	}

	PrintListNode(head)
	fmt.Println("")
	newHead := reverseKGroup(head, 3)
	PrintListNode(newHead)
}

func PrintListNode(head *ListNode) {
	cur := head
	for cur != nil {
		fmt.Print(cur.Val, " ")
		cur = cur.Next
	}
}
//leetcode 25

type ListNode struct {
	Val  int
	Next *ListNode
}

//快慢指针+递归
/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func reverseKGroup(head *ListNode, k int) *ListNode {
	test := head
	v := k
	for i := 0; i < k; i++ {
		if test != nil {
			test = test.Next
		} else {
			//没有足够的个数
			return head
		}
	}

	tail := head //现在的头是最后的尾巴
	//一般来说当前状态都是prev在前面，cur和next在当前处理节点
	var prev *ListNode
	cur := head
	next := head
	//迭代v次
	for i := 1; i <= v; i++ {
		next = cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}

	tail.Next = reverseKGroup(cur, k)
	return prev
}
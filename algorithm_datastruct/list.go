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

//leetcode 206
/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func reverseList(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}
	var prev *ListNode
	cur := head
	next := head
	for cur != nil {
		next = cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}
	return prev
}

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	fakeHead := &ListNode{}
	cur1 := l1
	cur2 := l2
	cur := fakeHead
	for cur1 != nil && cur2 != nil {
		if cur1.Val < cur2.Val {
			cur.Next = cur1
			cur1 = cur1.Next
		} else {
			cur.Next = cur2
			cur2 = cur2.Next
		}
		cur = cur.Next
	}
	if cur1 != nil {
		cur.Next = cur1
	}
	if cur2 != nil {
		cur.Next = cur2
	}
	return fakeHead.Next
}

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
//  判断成不成环，使用一个快慢指针，如果最终两个指针可以重合，表示有环
func hasCycle(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}
	faster := head.Next.Next
	slower := head.Next

	for faster != nil && slower != nil {
		if faster == slower {
			return true
		}
		if faster.Next != nil {
			faster = faster.Next.Next
		} else {
			faster = nil
		}
		slower = slower.Next
	}
	return false
}
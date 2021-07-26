package recursion

import "math"

type TreeNode struct {
	Val int
	Left *TreeNode
	Right *TreeNode
}

// leetcode 98
func isValidBST(root *TreeNode) bool {
	return _isValidBST(root, math.Inf(-1), math.Inf(1))
}

func _isValidBST(root *TreeNode, min float64, max float64) bool {
	if root == nil {
		return true
	}
	v := float64(root.Val)
	return min < v && v < max && _isValidBST(root.Left,min, v) && _isValidBST(root.Right, v, max)
}

// leetcode 99
func recoverTree(root *TreeNode)  {
	var first, second, pre *TreeNode
	inOrder(root, &first, &second, &pre)
	if first != nil && second != nil {
		first.Val, second.Val = second.Val, first.Val
	}
}

func inOrder(root *TreeNode, first, second, pre **TreeNode)  {
	if root == nil {
		return
	}

	inOrder(root.Left, first, second, pre)

	if *pre != nil && (*pre).Val >= root.Val {
		// 如果first还没有取值的话，那个pre的数字是错误的
		if *first == nil {
			*first = *pre
		}
		// second的是当前的节点是错误的，自己画个图就知道了
		*second = root
	}
	*pre = root
	inOrder(root.Right, first, second, pre)
}

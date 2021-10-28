package tree


type TreeNode struct {
     Val int
     Left *TreeNode
     Right *TreeNode
}

// leetcode
func sumNumbers(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return travel(root, 0)
}

// 使用递归做深度优先搜索，每层的和都要先乘以10
// 递归的终止条件是没有左右节点了，自己就是叶子节点，就把到自己的路径和加起来就可以了
// node 是遍历到的当前节点，cur是从root到node（不包含node）的路径的数字和
// node 的和有两种情况，一种是没有子节点，一种是有子节点；有子节点的话，是子节点的和，没有的话，就是到自己的这个路径和
func travel(node *TreeNode, cur int) int {
	// 终止条件
	cur = 10 * cur + node.Val
	if node.Left == nil && node.Right == nil {
		return cur
	}
	// 递归条件
	sum := 0
	if node.Left != nil {
		sum += travel(node.Left, cur)
	}
	if node.Right != nil {
		sum += travel(node.Right, cur)
	}
	return sum
}
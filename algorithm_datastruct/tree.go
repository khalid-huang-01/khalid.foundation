// 遍历
// leetcode 94
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
//  中序遍历
func inorderTraversal(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	result := make([]int, 0)
	inorder(root, &result)
	return result
}

func inorder(root *TreeNode, result *[]int) {
	if root == nil {
		return
	}
	inorder(root.Left, result)
	*result = append(*result, root.Val)
	inorder(root.Right, result)
}

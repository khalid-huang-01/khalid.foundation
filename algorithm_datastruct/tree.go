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

// 100. Same Tree
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func isSameTree(p *TreeNode, q *TreeNode) bool {
	if p == nil && q == nil {
		return true
	}
	if p == nil || q == nil {
		return false
	}
	//前序判断
	if p.Val == q.Val && isSameTree(p.Left, q.Left) && isSameTree(p.Right, q.Right) {
		return true
	}
	return false
}

// leetcode 404
// 本质还是遍历，就是添加了一点变化，前序遍历
func sumOfLeftLeaves(root *TreeNode) int {
	result := 0
	travel(root, &result, 0)
	return result
}

func travel(root *TreeNode, result *int, flag int) {
	if root == nil {
		return
	}
	if root.Left == nil && root.Right == nil && flag == -1 {
		*result = *result + root.Val
		return
	}
	travel(root.Left, result, -1)
	travel(root.Right, result, 1)
}

// leetcode 429
/**
 * Definition for a Node.
 * type Node struct {
 *     Val int
 *     Children []*Node
 * }
 */
// 在层搜索上，可以使用广度搜索+nil的方法，也可以使用记录层级的方式
//  广度搜索
func levelOrder(root *Node) [][]int {
	result := make([][]int, 0)
	if root == nil {
		return result
	}

	queue := make([]*Node, 0)
	queue = append(queue, root)
	queue = append(queue, nil)
	for len(queue) != 1 { //这里是1，因为可能最后只剩下一个nil
		tmp := make([]int, 0)
		// flag := false
		// fmt.Println(queue)
		for queue[0] != nil {
			cur := queue[0]
			tmp = append(tmp, cur.Val)
			queue = append(queue, cur.Children...)
			queue = queue[1:]
		}
		queue = queue[1:]
		queue = append(queue, nil)
		result = append(result, tmp)
	}
	return result
}
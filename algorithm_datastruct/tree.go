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
// 学习下https://leetcode-cn.com/problems/n-ary-tree-level-order-traversal/solution/ncha-shu-de-ceng-xu-bian-li-by-leetcode/
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

// leetcode 105
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
//  遍历前序建立根节点，利用中序建立关系
var pCurIdx int

func buildTree(preorder []int, inorder []int) *TreeNode {
	pCurIdx = 0
	return helper(preorder, inorder, 0, len(inorder)-1)
}

func helper(preorder []int, inorder []int, iStartIdx int, iEndIdx int) *TreeNode {
	if iStartIdx > iEndIdx {
		return nil
	}

	root := &TreeNode{Val: preorder[pCurIdx]}
	pCurIdx++

	//切分左树和右树
	var idx int
	for idx = iStartIdx; idx <= iEndIdx; idx++ {
		if inorder[idx] == root.Val {
			break
		}
	}

	root.Left = helper(preorder, inorder, iStartIdx, idx-1)
	root.Right = helper(preorder, inorder, idx+1, iEndIdx)

	return root
}

// 根据中序和后序重建树
// 核心在于后序的倒序就是中右左，与前序的中左右是差不多的
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
var pCurIdx int

func buildTree(inorder []int, postorder []int) *TreeNode {
	pCurIdx = len(postorder) - 1
	return helper(inorder, postorder, 0, len(inorder)-1)
}

func helper(inorder []int, postorder []int, iStartIdx int, iEndIdx int) *TreeNode {
	if iStartIdx > iEndIdx {
		return nil
	}
	//根节点
	root := &TreeNode{Val: postorder[pCurIdx]}
	pCurIdx--

	//查找root的位置
	var idx int
	for idx = iStartIdx; idx <= iEndIdx; idx++ {
		if inorder[idx] == root.Val {
			break
		}
	}
	//右树
	root.Right = helper(inorder, postorder, idx+1, iEndIdx)
	root.Left = helper(inorder, postorder, iStartIdx, idx-1)

	return root
}

//  tries 820
//  s
func minimumLengthEncoding(words []string) int {
	nodes := make(map[*TrieNode]int, 0) // 记录每个树叶子结点及其对应的words里面的下标
	// 输入每个words，构建前缀树
	// 根结点
	root := NewTrieNode()
	for index, word := range words {
		bytes := []byte(word)
		cur := root
		// 反向插入
		for i := len(bytes) - 1; i >= 0; i-- {
			cur = cur.getOrInsert(bytes[i])
		}
		// 记录单词和对应的下标，方便后面记录
		nodes[cur] = index
	}
	// 统计叶子中不会成为前缀的,也就是count为0
	var result int
	for key, value := range nodes {
		if key.count == 0 {
			result = result + len(words[value]) + 1
		}
	}
	return result
}

type TrieNode struct {
	children []*TrieNode // 为26个字母，如果为空的话，表示没有这个字母
	// isEnd bool
	count int // 用count来代表end，方便处理
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		count:    0,
		children: make([]*TrieNode, 26),
	}
}

func (t *TrieNode) getOrInsert(char byte) *TrieNode {
	if t.children[char-'a'] == nil {
		t.children[char-'a'] = NewTrieNode()
		t.count++
	}
	return t.children[char-'a']
}

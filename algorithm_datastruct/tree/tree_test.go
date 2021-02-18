package tree

import (
	"fmt"
	"testing"
)

// 遍历
// leetcode 94
// Definition for a binary tree node.

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

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
 */
type Node struct {
	Val      int
	Children []*Node
}

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
//var pCurIdx int

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

func buildTree1(inorder []int, postorder []int) *TreeNode {
	pCurIdx = len(postorder) - 1
	return helper1(inorder, postorder, 0, len(inorder)-1)
}

func helper1(inorder []int, postorder []int, iStartIdx int, iEndIdx int) *TreeNode {
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
	root.Right = helper1(inorder, postorder, idx+1, iEndIdx)
	root.Left = helper1(inorder, postorder, iStartIdx, idx-1)

	return root
}

//  tries 820
//func minimumLengthEncoding1(words []string) int {
//	nodes := make(map[*TrieNode]int, 0) // 记录每个树叶子结点及其对应的words里面的下标
//	// 输入每个words，构建前缀树
//	// 根结点
//	root := NewTrieNode()
//	for index, word := range words {
//		bytes := []byte(word)
//		cur := root
//		// 反向插入
//		for i := len(bytes) - 1; i >= 0; i-- {
//			cur = cur.getOrInsert(bytes[i])
//		}
//		// 记录单词和对应的下标，方便后面记录
//		nodes[cur] = index
//	}
//	// 统计叶子中不会成为前缀的,也就是count为0
//	var result int
//	for key, value := range nodes {
//		if key.count == 0 {
//			result = result + len(words[value]) + 1
//		}
//	}
//	return result
//}
//
//type TrieNode struct {
//	val      byte        // 该节点所代表的字符，val - 'a'其实就等于这个节点在父节点的儿子列表中的下标
//	children []*TrieNode // 为26个字母，如果为空的话，表示没有这个字母，每个下标对应1个字母，比如下标1就代表a
//	// isEnd bool
//	count int // 用count来代表end，方便处理，count == len(children)，当count为零时，表示是最后一个字母
//}
//
//func NewTrieNode() *TrieNode {
//	return &TrieNode{
//		count:    0,
//		children: make([]*TrieNode, 26),
//	}
//}
//
//func (t *TrieNode) getOrInsert(char byte) *TrieNode {
//	if t.children[char-'a'] == nil {
//		t.children[char-'a'] = NewTrieNode()
//		t.count++
//	}
//	return t.children[char-'a']
//}
type TrieNode struct {
	children []*TrieNode // 为26个字母，如果为空的话，表示没有这个字母，每个下标对应1个字母，比如下标1就代表a(通过letter - 'a'得到)
	childCount int
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children:   make([]*TrieNode, 26),
		childCount: 0, // 当childCount为0的时候，表示这是一个结尾的字母，没有后续的字母了
	}
}

func (t *TrieNode) getOrInsert(letter byte) *TrieNode {
	if t.children[letter - 'a'] == nil {
		t.children[letter - 'a'] = NewTrieNode()
		t.childCount++
	}
	return t.children[letter - 'a']
}

func minimumLengthEncoding(words []string) int {
	// 建立 前缀树
	root := NewTrieNode()
	nodes := make(map[*TrieNode]int, 0) // 记录每个单词下标及其叶子节点（最后一个单词）的对应关系
	for index, word := range words {
		cur := root
		letters := []byte(word)
		for i := len(letters) - 1; i >= 0; i-- {
			cur = cur.getOrInsert(letters[i])
		}
		nodes[cur] = index
	}
	// 获取答案
	rsl := 0
	for node, index := range nodes {
		if node.childCount == 0 {
			rsl += len(words[index]) + 1
		}
	}
	return rsl
}

// leetcode 208
type Trie struct {
	childrens []*Trie // 26个字母，有值表示有这个儿子，没有的话就没有
	isEnd bool // 代表是不是一个单词的结束
}


/** Initialize your data structure here. */
func Constructor() Trie {
	return Trie{
		childrens: make([]*Trie, 26),
		isEnd:     false,
	}
}

func (this *Trie) getOrInsert(letter byte) *Trie {
	if this.childrens[letter-'a'] == nil{
		this.childrens[letter - 'a'] = &Trie{
			childrens: make([]*Trie, 26),
			isEnd:     false,
		}
	}
	return this.childrens[letter- 'a']
}
 
/** Inserts a word into the trie. */
func (this *Trie) Insert(word string)  {
	letters := []byte(word)
	cur := this
	for _, letter := range letters {
		cur = cur.getOrInsert(letter)
	}
	// 给最后添加上单词标记
	cur.isEnd = true
}



/** Returns if the word is in the trie. */
func (this *Trie) Search(word string) bool {
	letters := []byte(word)
	cur := this
	for _, letter := range letters {
		if cur.childrens[letter-'a'] == nil {
			return false
		}
		cur = cur.childrens[letter-'a']
	}
	// 判断最后一个字符是不是结束字符
	return cur.isEnd
}


/** Returns if there is any word in the trie that starts with the given prefix. */
func (this *Trie) StartsWith(prefix string) bool {
	letters := []byte(prefix)
	cur := this
	for _, letter := range letters {
		if cur.childrens[letter-'a'] == nil {
			return false
		}
		cur = cur.childrens[letter-'a']
	}
	return true
}

func TestTrie(t *testing.T) {
	obj := Constructor()
	word := "apple"
	obj.Insert(word)
	fmt.Println(obj.Search(word))
	fmt.Println(obj.StartsWith("app"))
}


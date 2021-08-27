package stack

// leetcode 173
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type BSTIterator struct {
	stack []*TreeNode
}

func Constructor(root *TreeNode) BSTIterator {
	bi := BSTIterator{[]*TreeNode{}}
	bi.push(root)
	return bi
}

func (this *BSTIterator) push(node *TreeNode) {
	for node != nil {
		this.stack = append(this.stack, node)
		node = node.Left // 一直往左边找，同时把当前节点给保持起来
	}
}

func (this *BSTIterator) Next() int {
	tmp := this.stack[len(this.stack)-1]
	this.stack = this.stack[:len(this.stack)-1]
	this.push(tmp.Right) // 把右边的给入栈
	return tmp.Val
}

func (this *BSTIterator) HasNext() bool {
	return len(this.stack) > 0
}

// leetcode 20
func isValid(s string) bool {
	stack := make([]rune, 0)
	pairs := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}
	for _, v := range s {
		switch v {
		case '(', '{', '[':
			stack = append(stack, v)
		case ')', '}', ']':
			if len(stack) == 0 || stack[len(stack)-1] != pairs[v] {
				return false
			}
			stack = stack[:len(stack)-1]
		default:
			return false
		}
	}
	return len(stack) == 0
}

// leetcode 28
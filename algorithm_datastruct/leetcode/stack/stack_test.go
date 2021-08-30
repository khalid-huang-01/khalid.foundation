package stack

import "testing"

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

func nextExceed(array []int) []int {
	monoStack := make([]int, 0)
	rsl := make([]int, len(array))
	for i, v := range array {
		for len(monoStack) != 0 && array[monoStack[len(monoStack)-1]] < v {
			rsl[monoStack[len(monoStack)-1]] = i - monoStack[len(monoStack)-1]
			monoStack = monoStack[:len(monoStack)-1]
		}
		monoStack = append(monoStack, i)
	}
	for i := 0; i < len(monoStack); i++ {
		rsl[monoStack[i]] = -1
	}
	return rsl
}

func TestNextExceed(t *testing.T) {
	array := []int{5,3,1,2,4}
	t.Log(nextExceed(array))
}

// leetcode 84
func largestRectangleArea(heights []int) int {
	monoStack := make([]int, 0)
	result := 0
	heights = append(heights, 0) // 放入一个0，保证全部的都可以计算到，特别是以1的顶为上边框的
	for i, v := range heights {
		// 进一个元素，要弹出可以由这个元素得到结果的前面的元素
		for len(monoStack) != 0 && heights[monoStack[len(monoStack)-1]] > v {
			top := len(monoStack) - 1
			// 确定高度
			height := heights[monoStack[top]]
			// 优化：不断弹出相同元素
			for len(monoStack) != 0 && heights[monoStack[len(monoStack)-1]] == height {
				monoStack = monoStack[:len(monoStack) - 1]
			}
			// 计算长度，分两种情况，一种是栈里没有元素了，表明这个栈里面没有比它要小的左边界，那么它的宽度就是从0到自己的下标位置； 另
			// 另一个种是栈里面有比它小的左边界，也就是当前的栈顶元素，因为比它大的都被弹出去了
			var width int
			if len(monoStack) == 0 {
				width = i
			} else {
				width = i - monoStack[len(monoStack)-1] - 1
			}
			// 计算面积
			result = max(result, width * height)
		}
		monoStack = append(monoStack, i)
	}
	return result
}

func TestLargestRectangleArea(t *testing.T) {
	//array := []int{2,1,5,6,2,3} // 答案是10
	array := []int{6,5,6} // 答案是15
	result := largestRectangleArea(array)
	t.Log(result)
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a

}





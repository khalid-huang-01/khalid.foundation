package main

import (
	"fmt"
	"testing"
)

// leetcode 739
func dailyTemperatures(T []int) []int {
	stack := make([]int, 0)
	result := make([]int, len(T))
	for i := 0; i < len(T); i++ {
		for len(stack) != 0 && T[stack[len(stack)-1]] < T[i] {
			result[stack[len(stack)-1]] = i - stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}
	return result
}

// 《算法竞赛入门经典》 6.1.2 铁轨
func solution_6_1_2(num int, target []int) bool {
	stack := make([]int, 0)
	A, targetIndex := 1, 0 // A是原始车编号（入栈的） B是目标车编号下标
	// 在目标车编号没有结束前
	for targetIndex < num {
		// 要进站的车编号就是要找的
		if A == target[targetIndex] {
			A++
			targetIndex++
		} else if len(stack) != 0 && stack[len(stack)-1] == target[targetIndex] {
			//栈里面有，之前添加过
			stack = stack[:len(stack)-1] //直接出栈
			targetIndex++
		} else if A <= num {
			// 把一个新的入栈
			stack = append(stack, A)
			A++
		} else {
			// 上面的都不行，直接false
			return false
		}
	}
	return true
}

// https://leetcode-cn.com/problems/maximum-frequency-stack/
// 算法的思想就是把用两个map来保存信息，构建过程如下，依次输入 5, 7,5, 7, 4, 5
// 输入5 : val2Freq(5:1) freq2Vals(1:[5]) //表示当前5这个数字出现频率为1，频率为1的数字有5 // 我们用这个数字来保存了5的历史，这样当把频率为2的5打出去的时候，其实我们这个map还保留了频率为1的5的情况
// 输入7： val2Freq(5:1, 7:1) freq2Vals(1:[5,7]) //出栈的顺序其实就是里面列表的反序了
// 输入5： val2Freq(5:2, 7:1) freq2Vals(1:[5,7], 2:[5]) // 这样把5弹出去的时候，就需要恢复到上面的一步
type FreqStack struct {
	val2Freq  map[int]int
	freq2Vals map[int][]int
	maxFreq   int
}

func Constructor() FreqStack {
	f := FreqStack{
		val2Freq:  make(map[int]int, 0), // 用于知道当前添加进来的数字在栈里有几个，好放到合适的freq2Vals里面
		freq2Vals: make(map[int][]int, 0),
		maxFreq:   0,
	}
	return f
}

func (f *FreqStack) Push(x int) {
	f.val2Freq[x] += 1
	if f.freq2Vals[f.val2Freq[x]] == nil {
		f.freq2Vals[f.val2Freq[x]] = make([]int, 0)
	}
	f.freq2Vals[f.val2Freq[x]] = append(f.freq2Vals[f.val2Freq[x]], x)
	// 判断当前最大的频率是多少
	if f.maxFreq < f.val2Freq[x] {
		f.maxFreq = f.val2Freq[x]
	}
}

func (f *FreqStack) Pop() int {
	size := len(f.freq2Vals[f.maxFreq])
	result := f.freq2Vals[f.maxFreq][size-1]
	f.freq2Vals[f.maxFreq] = f.freq2Vals[f.maxFreq][:size-1] // 消除要弹出去的元素
	if len(f.freq2Vals[f.maxFreq]) == 0 {
		f.maxFreq -= 1 // 当前最大的已经没有了，往第二频率梯队找
	}
	return result
}



/*
题目描述：
给一个数组，返回一个大小相同的数组。返回的数组的第i个位置的值应当是，对于原数组中的第i个元素，至少往右走多少步，才能遇到一个比自己大的元素（如果之后没有比自己大的元素，或者已经是最后一个元素，则在返回数组的对应位置放上-1）。
简单的例子：
input: 5,3,1,2,4
return: -1 3 1 1 -1
 */
func nexExceed(array []int) []int {
	result := make([]int, len(array))
	monoStack := make([]int,0)
	for i, v := range array {
		// monoStack[len(monoStack-1)] 表示栈顶元素
		for len(monoStack) != 0 && (array[monoStack[len(monoStack)-1]] < v ) {
			result[monoStack[len(monoStack)-1]] = i - monoStack[len(monoStack)-1]
			monoStack = monoStack[:len(monoStack)-1]
		}
		monoStack = append(monoStack, i)
	}
	for len(monoStack) != 0 {
		result[monoStack[len(monoStack)-1]] = -1
		monoStack = monoStack[:len(monoStack)-1]
	}
	return result
}

func TestNextExceed(t *testing.T) {
	array := []int{5,3,1,2,4}
	result := nexExceed(array) // 结果应该是-1，3，1，1，-1
	t.Log(result)
}

//-------------------
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
			// 不断弹出相同元素
			for len(monoStack) != 0 && heights[monoStack[len(monoStack)-1]] == height {
				monoStack = monoStack[:len(monoStack) - 1]
			}
			// 计算长度，分两种情况，一种是栈里没有元素了，表明这个栈里面没有比它要小的左边界，那么它的宽度就是从0到自己的下标位置； 另
			// 另一个种是栈里面有比它小的左边界，也就是栈顶的元素，因为比它大的都被弹出去了
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

func main() {
	num := 6
	target := []int{6, 5, 4, 3, 2, 1}
	fmt.Println(solution_6_1_2(num, target))
}

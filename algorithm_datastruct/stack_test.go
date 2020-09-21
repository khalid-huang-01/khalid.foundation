package main

import "fmt"

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

func main() {
	num := 6
	target := []int{6, 5, 4, 3, 2, 1}
	fmt.Println(solution_6_1_2(num, target))
}

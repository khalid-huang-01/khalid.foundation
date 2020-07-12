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

func main() {
	num := 6
	target := []int{6,5,4,3,2,1}
	fmt.Println(solution_6_1_2(num, target))
}
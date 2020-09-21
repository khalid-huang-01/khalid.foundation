package main

import (
	"fmt"
	"testing"
)

// nums 为列表，可正数，可负数，长度一定大于2
// 设置f(n)表示从下标0到下标n的最大不相邻数列之和，那么就可以有f(n) = max(f(n-1), f(n-2) + v(n))，其中v(n) = max(下标n的值，0)
func maxNotAdjacentNumsSum(nums []int) int {
	f := make([]int, len(nums))
	f[0] = max(0, nums[0])
	f[1] = max(f[0], nums[1])
	for i := 2; i < len(nums); i++ {
		v := max(0, nums[i])
		f[i] = max(f[i-1], f[i-2]+v)
	}
	return f[len(nums)-1]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// n表示游戏的数目，x表示总天数
// values表示通关游戏所获得的成就值，cost表示需要花费的天数
func zeroOnePackage(n int, x int, values []int, costs []int) int {
	f := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		f[i] = make([]int, x+1)
	}
	for i := 1; i <= n; i++ {
		for j := 0; j <= x; j++ {
			// 先初始化(假定不使用)，如果是只有一款就
			if i == 1 {
				f[i][j] = 0
			} else {
				f[i][j] = f[i-1][j]
			}
			// 比较使用的情况，要满足剩余的天数大于消耗的天数
			if j >= costs[i] && f[i][j] < f[i-1][j-costs[i]]+values[i] {
				f[i][j] = f[i-1][j-costs[i]] + values[i]
			}
		}
	}
	return f[n][x]
}

func TestZeroOnePackage(t *testing.T) {
	//nums := []int{7,8,-1,5,6,10,-3}
	//fmt.Println(maxNotAdjacentNumsSum(nums))

	n := 2
	x := 2
	values := []int{0, 10, 20}
	costs := []int{0, 1, 2}
	fmt.Println(zeroOnePackage(n, x, values, costs))
}

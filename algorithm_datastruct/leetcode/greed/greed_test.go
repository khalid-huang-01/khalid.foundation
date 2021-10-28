package greed

import "sort"

// leetcode 1833
func maxIceCream(costs []int, coins int) int {
	sort.Ints(costs)
	cur, rsl := 0, 0
	for i := 0; i < len(costs); i++ {
		cur += costs[i]
		if cur > coins {
			break
		}
		rsl += 1
	}
	return rsl
}

// leetcode 1558
// 使用贪心策略：每轮都尽量把所有的数字除2，如果除不尽，就加一次减一操作
func minOperations(nums []int) int {
	allZero := false
	count := 0
	for allZero == false {
		allZero = true
		for i := 0; i < len(nums); i++ {
			if nums[i] % 2 == 1 {
				count += 1 //op 1
			}
			nums[i] /= 2
			if nums[i] != 0 {
				allZero = false
			}
		}
		count += 1 // op 2
	}
	return count - 1 // 最后的（0,1）或者（0，0）会多算一次
}
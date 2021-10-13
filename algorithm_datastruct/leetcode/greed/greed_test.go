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


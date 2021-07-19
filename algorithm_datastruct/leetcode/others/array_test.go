package others

import (
	"testing"
)

// leetcode 189
func rotate(nums []int, k int)  {
	// 可能绕几圈
	k %= len(nums)
	reverseArray(nums, 0, len(nums)-1)
	reverseArray(nums, 0, k-1)
	reverseArray(nums, k, len(nums)-1)
}

func reverseArray(nums []int, left, right int) {
	for left < right {
		nums[left], nums[right] = nums[right], nums[left]
		left++
		right--
	}
}


// -----------------------------
// {2, 4, 1, 16, 7, 5, 11, 9} -> 11
func maxDiff(nums []int) int {
	dp := make([]int, len(nums))
	dp[0] = nums[0]
	for i := 1; i < len(nums); i++ {
		dp[i] = max(dp[i-1], nums[i])
	}
	rsl := 0
	for i := 0; i < len(nums); i++ {
		diff := dp[i] - nums[i]
		rsl = max(rsl, diff)
	}
	return rsl
}

func TestMaxDiff(t *testing.T) {
	nums := []int{2, 4, 1, 16, 7, 5, 11, 9}
	t.Log(maxDiff(nums))
}



func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
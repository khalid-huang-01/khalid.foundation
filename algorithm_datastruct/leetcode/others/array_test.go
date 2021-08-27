package others

import (
	"testing"
	"time"
)

// leetcode 350
// 无序的且重复数字也要多次返回
// 使用map，一个数组计数，一个消除，消除就加入
func intersect(nums1 []int, nums2 []int) []int {
	counts := make(map[int]int)
	for i := 0; i < len(nums1); i++ {
		counts[nums1[i]]++
	}
	rsl := make([]int, 0)
	for i := 0; i < len(nums2); i++ {
		if counts[nums2[i]] > 0 {
			rsl = append(rsl, nums2[i])
			counts[nums2[i]]--
		}
	}
	return rsl
}

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
	t.Log("start")
	var intervalTime = 5 * time.Second
	time.Sleep(intervalTime)
	t.Log("end")
	nums := []int{2, 4, 1, 16, 7, 5, 11, 9}
	t.Log(maxDiff(nums))
}



func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

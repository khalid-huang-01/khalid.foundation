package slide_windows

import "math"

// leetcode 209
func minSubArrayLen(target int, nums []int) int {
	left, right, curSum := 0, 0, 0
	ans := math.MaxInt64
	for right < len(nums) {
		// 添加右侧数据
		curSum += nums[right]

		// 通过移动左侧数据来维护条件，同时获取最新结果
		for curSum >= target {
			// 判断是不是比现有的结果小，如果是的化，更新
			if ans > right - left + 1 {
				ans = right - left + 1
			}
			curSum -= nums[left]
			left = left + 1
		}
		// 移动右侧
		right += 1
	}
	if ans == math.MaxInt64 {
		ans = 0
	}
	return ans
}

// leedcode 3
func lengthOfLongestSubstring(s string) int {
	// 统计字母个数
	count := make(map[byte]int)
	left, right := 0, 0
	ans := 0
	for right < len(s) {
		// 加入右边数字
		count[s[right]]++

		// 调整左侧，左侧数字往右移动，使条件满足;并获取最新结果
		for count[s[right]] != 1 {
			count[s[left]] -= 1
			left += 1
		}
		if ans < right - left + 1 {
			ans = right - left + 1
		}

		// 移动右侧
		right += 1
	}
	return ans
}
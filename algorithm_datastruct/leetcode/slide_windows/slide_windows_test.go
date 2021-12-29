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

// leetcode 1004
func longestOnes(nums []int, k int) int {
	left, right := 0, 0
	cost, rsl := 0, 0
	for right < len(nums) {
		// 加入右侧
		if nums[right] == 0 {
			cost += 1 // 变更一次
		}
		
		// 通过右移左侧游标，使条件满足，主要就是不让变更的数量大于k
		for cost > k {
			if nums[left] == 0 {
				cost -= 1 // 撤销前面变更的
			}
			left += 1
		}
		
		// 获取最新结果
		if rsl < right - left + 1 {
			rsl = right - left + 1
		}

		// 推进右侧游标
		right += 1
	}
	return rsl
}

// leetcode 1208
func abs(a, b byte) int {
	tmp := int(a) - int(b)
	if tmp <  0 {
		return -tmp
	}
	return tmp
}

func equalSubstring(s string, t string, maxCost int) int {
	left, right := 0, 0
	cost, rsl := 0, 0
	for right < len(s) {
		// 加入右侧数据，比较
		cost += abs(s[right], t[right])

		// 右移左侧游标，满足条件
		for cost > maxCost {
			cost -= abs(s[left], t[left])
			left += 1
		}

		// 获取最新结果
		if rsl < right - left + 1 {
			rsl = right - left + 1
		}

		// 右移右侧游标
		right += 1
	}
	return rsl
}

// 双端队列解决
func maxSlidingWindow(nums []int, k int) []int {
	var result []int
	var window []int // 保存下标

	for index, value := range nums {
		// 维护双端队列的属性
		// 1. 判断队头是否超过区间（从队头出列）
		if index >= k && index-window[0] == k {
			window = window[1:]
		}
		// 2. 判断即将入队元素的位置，把比当前数字小的都从尾部出队，通过出队其他元素的方式（从队尾出队）
		for len(window) > 0 && nums[window[len(window)-1]] < value {
			window = window[:len(window)-1]
		}
		//入队
		window = append(window, index)
		//获取当前位置的结果
		if index >= k-1 {
			result = append(result, nums[window[0]])
		}
	}
	return result
}
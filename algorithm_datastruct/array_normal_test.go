package main

import (
	"strconv"
	"testing"
)

// leetcode 228
func summaryRanges(nums []int) []string {
	result := make([]string, 0)
	if len(nums) == 0 {
		return result
	}

	last := nums[0]
	count := 0
	for _, val := range nums[1:] {
		if val == last+count+1 {
			count++
			continue
		}
		if count == 0 {
			result = append(result, strconv.Itoa(last))
		} else {
			result = append(result, strconv.Itoa(last)+"->"+strconv.Itoa(last+count))
		}
		last = val
		count = 0
	}
	return result
}

// 剑指offer p82 面试题11
// 返回旋转数组中的最小值
// 比如 3 4 5 1 2，返回1
// 其实就是利用二分查找，每次都找到那个不是顺序排列的，因为最小值就在两个非顺序排序的交界处
// 如果说出现类似 1 0 1 1 1 这种的有大量重复的，就只能使用顺序查找的方式
func MinValueOfRotaryArray(nums []int) int {
	low := 0
	high := len(nums) - 1

	var mid int
	for low <= high {
		if high-low == 1 {
			mid = high // 如果相关为1，证明应该是在右边的
			break
		}
		// 寻找交界处
		mid = low + (high-low)/2
		if nums[mid] >= nums[low] {
			low = mid
		} else if nums[mid] < nums[high] {
			high = mid
		}
	}
	return nums[mid]
}

func TestMinValueOfRotaryArray(t *testing.T) {
	nums := []int{3, 4, 5, 0, 1, 2}
	t.Log("value: ", MinValueOfRotaryArray(nums))
}

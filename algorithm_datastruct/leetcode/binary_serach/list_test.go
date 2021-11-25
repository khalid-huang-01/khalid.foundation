package binary_serach

import (
	"testing"
)

func _binarySearch(array []int, target int, left int, right int) int {
	// 终止条件
	if left > right {
		return -1
	}
	mid := left + (right-left)/2
	if array[mid] == target {
		return mid
	}
	if array[mid] < target {
		return _binarySearch(array, target, mid+1, right)
	}
	return _binarySearch(array, target, left, mid-1)
}

// 二分查找
func binarySearch(array []int, target int) int {
	left, right := 0, len(array)-1
	return _binarySearch(array, target, left, right)
}

func binarySearch2(array []int, target int) int {
	low, high := 0, len(array)-1
	rsl := -1
	for low <= high {
		mid := low + (high-low)/2
		if array[mid] == target {
			rsl = mid
			break
		}
		if array[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return rsl
}

func TestBinarySearch(t *testing.T) {
	// 公用的配置
	t.Run("success", func(t *testing.T) {
		array := []int{3, 4, 15, 17, 89, 167, 345}
		target := 4
		want := 1
		rsl := binarySearch2(array, target)
		if want != rsl {
			t.Errorf("got %d, want %d", rsl, want)
		} else {
			t.Log("success")
		}
	})
	t.Run("fail", func(t *testing.T) {
		array := []int{3, 4, 15, 17, 89, 167, 345}
		target := 5
		want := -1
		rsl := binarySearch2(array, target)
		if want != rsl {
			t.Errorf("got %d, want %d", rsl, want)
		} else {
			t.Log("success")
		}
	})
}

// leetcode 34
func searchLeftBound(nums []int, target int, low int, high int) int {
	for low <= high {
		mid := low + (high-low)/2
		// 确定是左边界
		if nums[mid] == target && (mid == 0 || nums[mid-1] != nums[mid]) {
			return mid
		}
		if nums[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}

func searchRightBound(nums []int, target int, low int, high int) int {
	for low <= high {
		mid := low + (high-low)/2
		if nums[mid] == target && (mid == len(nums)-1 || nums[mid+1] != nums[mid]) {
			return mid
		}
		// 确定下如果实在连续目标等值里面要怎么做，这里是要往上面走, 剩下的就是如果taget再右侧，就把low往上提
		if nums[mid] <= target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}

// leetcode 34
func searchRange(nums []int, target int) []int {
	firstPosition := searchLeftBound(nums, target, 0, len(nums)-1)
	lastPostion := searchRightBound(nums, target, 0, len(nums)-1)
	return []int{firstPosition, lastPostion}
}

// leetcode 875
func minEatingSpeed(piles []int, h int) int {
	low := 1
	high := maxPiles(piles)
	var count int
	for low < high {
		mid = low + (high-low)/2
		count = calTimes(piles, mid)
		if count > h {
			low = mid + 1
		} else {
			high = mid
		}
	}
	return low
}

func calTimes(piles []int, speed int) int {
	count := 0
	speedFloat := float64(speed)
	for _, value := range piels {
		count += int(math.Ceil(float64(value) / speedFloat))
	}
	return count
}

func maxPiles(piles []int) int {
	value := piles[0]
	for i := 1; i < len(piles); i++ {
		if value < piles[i] {
			value = piles[i]
		}
	}
	return value
}

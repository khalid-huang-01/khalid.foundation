// 二分查找代码及其变种

package main

import (
	"fmt"
)

func binarySearchWithRecursion(arr []int, target int, low int, high int) int {
	if low > high {
		return -1
	}
	middle := low + (high-low)/2
	if target == arr[middle] {
		return middle
	}
	if target < arr[middle] {
		return binarySearchWithRecursion(arr, target, low, middle-1)
	} else {
		return binarySearchWithRecursion(arr, target, middle+1, high)
	}
}

func binarySearchWithoutRecursion(arr []int, target int) int {
	low := 0
	high := len(arr)
	for low <= high {
		middle := low + (high-low)/2
		if target == arr[middle] {
			return middle
		}
		if target < arr[middle] {
			high = middle - 1
		} else {
			low = middle + 1
		}
	}

	return -1
}

func testWithRecursion() {
	arr := []int{2, 3, 6, 8, 9, 11, 15, 56}
	target := 12
	result := binarySearchWithRecursion(arr, target, 0, len(arr)-1)
	fmt.Println("result: ", result)
	return
}

func testWithoutRecursion() {
	arr := []int{2, 3, 6, 8, 9, 11, 15, 56}
	target := 12
	result := binarySearchWithoutRecursion(arr, target)
	fmt.Println("result: ", result)
	return
}

func searchLowBoundary(nums []int, target int) int {
	low := 0
	high := len(nums)
	for low <= high {
		middle := low + (high-low)/2
		if target == nums[middle] && (middle == 0 || nums[middle-1] != nums[middle]) {
			return middle
		}
		if target <= nums[middle] {
			high = middle - 1
		} else {
			low = middle + 1
		}
	}
	return -1
}

func searchHighBoundary(nums []int, target int) int {
	low := 0
	high := len(nums)
	for low <= high {
		middle := low + (high-low)/2
		if target == nums[middle] && (middle == len(nums)-1 || nums[middle+1] != nums[middle]) {
			return middle
		}
		if target < nums[middle] {
			high = middle - 1
		} else {
			low = middle + 1
		}
	}
	return -1
}

func searchRange(nums []int, target int) []int {
	low := searchLowBoundary(nums, target)
	high := searchHighBoundary(nums, target)
	return []int{low, high}
}

func leetcode_34() {
	arr := []int{5, 7, 7, 8, 8, 10}
	target := 8
	fmt.Println(searchRange(arr, target))
}

// 以对角线元素start为起点进行垂直或者平行的二分查找
func binarySearch(matrix [][]int, start int, target int, vertical bool) bool {
	var low, mid, high int
	low = start
	if vertical {
		high = len(matrix[0]) - 1
	} else {
		high = len(matrix) - 1
	}
	for low <= high {
		mid = low + (high - low) / 2
		// 列查找
		if vertical {
			if matrix[start][mid] == target {
				return true
			}
			if matrix[start][mid] < target {
				low = mid + 1
			} else {
				high = mid - 1
			}
		} else {
			// 行查找
			if matrix[mid][start] == target {
				return true
			}
			if matrix[mid][start] < target {
				low = mid + 1
			} else {
				high = mid - 1
			}
		}
	}
	return false
}

func searchMatrix(matrix [][]int, target int) bool {
	if matrix == nil || len(matrix) == 0 {
		return false
	}
	var shortDim int
	if len(matrix) < len(matrix[0]) {
		shortDim = len(matrix)
	} else {
		shortDim = len(matrix[0])
	}

	var verticalFound, horizontalFound bool
	var rsl bool
	for i := 0; i < shortDim; i++ {
		verticalFound = binarySearch(matrix, i, target, true)
		horizontalFound = binarySearch(matrix,i, target, false)
		if verticalFound || horizontalFound {
			rsl = true
			break
		}
	}
	return rsl
}

// leetcode 33
func search(nums []int, target int) int {
	var low, mid, high int
	low = 0
	high = len(nums) - 1
	for low <= high {
		mid = low + (high - low) / 2
		if nums[mid] == target {
			return mid
		}
		// 左边有序，如果有序满足情况，就在有序中找，否则在无序中找
		if nums[low] <= nums[mid] {
			if nums[low] <= target && target < nums[mid] {
				high = mid - 1
			} else {
				low = mid + 1
			}
		} else {
			// 右边有序
			if nums[mid] < target && target <= nums[high] {
				low = mid + 1
			} else {
				high = mid - 1
			}
		}
	}
	return -1
}

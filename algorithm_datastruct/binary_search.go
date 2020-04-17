// 二分查找代码及其变种

package main

import "fmt"

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
	for low < high {
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

func searchHighBoundary() {
	low := 0
	high := len(nums)
	for low < high {
		middle := low + (high-low)/2
		if target == nums[middle] && (middle == len(nums) || nums[middle+1] != nums[middle]) {
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

func main() {
	// testWithoutRecursion()
	// testWithRecursion()
	leetcode_34()
}

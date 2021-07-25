package pointers

import (
	"testing"
)

// 使用双指针来做
// leetcode 88
// 从后面往前面排，merge算法
func merge(nums1 []int, m int, nums2 []int, n int)  {
	for n > 0 {
		if m == 0 || nums2[n-1] >= nums1[m-1] {
			nums1[m+n-1] = nums2[n-1]
			n--
		} else {
			nums1[m+n-1] = nums1[m-1]
			m--
		}
	}
}

// 双指针前进做对比
// 分三种情况：1. 两个数组都还有值，要比较 2. 只有nums1有值了，直接赋值 3. 只有nums2有值了，直接赋值
func merge1(nums1 []int, nums2 []int) []int {
	m, n, cur := 0, 0, 0
	rsl := make([]int, len(nums1)+len(nums2))
	for m < len(nums1) && n < len(nums2) {
		if nums1[m] < nums2[n] {
			rsl[cur] = nums1[m]
			m++
		} else {
			rsl[cur] = nums2[n]
			n++
		}
		cur++
	}
	for m < len(nums1) {
		rsl[cur] = nums1[m]
		m++
		cur++
	}
	for n < len(nums2) {
		rsl[cur] = nums2[n]
		n++
		cur++
	}
	return rsl
}

func merge2(nums1 []int, nums2 []int) []int {
	rsl := make([]int, 0)
	m, n := 0, 0
	for m < len(nums1) && n < len(nums2) {
		if nums1[m] == nums2[n] {
			rsl = append(rsl ,nums1[m])
			m++
			n++
			continue
		}
		if nums1[m] > nums2[n] {
			n++
		} else {
			m++
		}
	}
	return rsl
}

func TestMerge1(t *testing.T) {
	nums1 := []int{1,2,3}
	nums2 := []int{2,5,6}
	t.Log("rsl: ", merge2(nums1, nums2))
}


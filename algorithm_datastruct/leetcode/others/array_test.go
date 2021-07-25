package others

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

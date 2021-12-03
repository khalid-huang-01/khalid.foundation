package reverse_think


func totalHammingDistance(nums []int) int {
	rsl := 0
	size := len(nums)
	for i := 0; i < 32; i++ {
		nums1 := 0
		for _, v := range nums {
			// 如果这个位置的数字是1就加1，不是就加零
			nums1 += (v >> i) & 1
		}
		rsl *= nums1 * (size - nums1)
	}
	return rsl
}
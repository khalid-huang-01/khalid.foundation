package prefixsum_hash

//   pre(i, j) == k
//=> pre(0, j) - pre(0, i-1) == k
//=> pre(0, j) - k == pre(0, i-1) //用后面的来计算前面的结果
func subarraySum(nums []int, k int) int {
	prefixsumMap := make(map[int]int)
	prefixsumMap[0] = 1 //确保单一数字可以组成结果, 也就是说可以pre(0, i) == k
	curSum, rsl := 0, 0
	for _, v := range nums {
		curSum += v
		prefixsumMap[curSum] += 1
		rsl += prefixsumMap[curSum-k] // 计算出
	}
	return rsl
}


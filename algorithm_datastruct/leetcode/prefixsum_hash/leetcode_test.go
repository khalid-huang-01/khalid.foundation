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

// leetcode 523
// pre(i, j) = k * n
// => pre(i,j) % k = (k*n) % k
// => pre(i,j) % k = 0
// => (pre(0,j) - pre(0, i-1)) % k = 0
// => pre(0,j) % k = pre(0, i-1) %k
// 条件是 j != i => 组合的数必须大于2, 而且要满足可以从小标为0开始的组合
func checkSubarraySum(nums []int, k int) bool {
	prefixsumMap := make(map[int]int)
	prefixsumMap[0] = -1 // 虚拟设置下表位置-1的值为0，这样才可以从下表为0开始的地方开始
	isExist := false
	curSum := 0
	for i, v := range nums {
		curSum += v
		tmp := curSum % k
		_, ok := prefixsumMap[tmp]
		if ok && prefixsumMap[tmp]+1 != i{ // 保证组合的数目必须大于2
			isExist = true
			break
		}
		if !ok {
			prefixsumMap[tmp] = i
		}
	}
	return isExist
}


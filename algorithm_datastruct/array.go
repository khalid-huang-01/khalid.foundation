// subarray
// leetcode 974
// 暴力直接找, 接近O(n^2)
func check(A []int, start int, K int, result *int) {
	size := len(A)
	var sum int
	for i := start; i < size; i++ {
		sum += A[i]
		if sum%K == 0 {
			result++
		}
	}
}

func subarraysDivByK(A []int, K int) int {
	var result int

	for i := 0; i < len(A); i++ {
		check(A, i, K, &result)
	}
	return result
}

//通过规律, O(n)
// % is remainder operator in C++ (and not a proper modulus). To get a positive number, we have to add by the base: a mod b = ((a % b) + b) % b.
// Look for anytime S[j] % k == S[i-1] % k and you know everything in between [i...j] must be divisible by K!!
// 上面的公式无法适用于i == 0的时候，所以需要额外需要判断有s[j] % K == 0，但是为了统一算法，我们可以定下s[-1] = 0，这时是可以符合条件的
// https://leetcode.com/problems/subarray-sums-divisible-by-k/discuss/584722/C%2B%2B-O(N)-Explained
func subarraysDivByK(A []int, K int) int {
	remainMap := make(map[int]int, 0)
	var sum, remain, result int
	remainMap[0] = 1 //s[-1] = 0
	for i := 0; i < len(A); i++ {
		sum += A[i]
		remain = ((sum % K) + K) % K
		if val, ok := remainMap[remain]; ok {
			result += val
		}
		remainMap[remain] = remainMap[remain] + 1
	}
	return result
}

// 在数组的处理中，一般都需要用于map来缩小运算时间
func majorityElement(nums []int) int {
	count := make(map[int]int, 0)
	majority := len(nums) / 2
	for _, val := range nums {
		count[val] = count[val] + 1
		if count[val] > n {
			return v
		}
	}
	return -1
}
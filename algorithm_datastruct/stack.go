// leetcode 739
func dailyTemperatures(T []int) []int {
	stack := make([]int, 0)
	result := make([]int, len(T))
	for i := 0; i < len(T); i++ {
		for len(stack) != 0 && T[stack[len(stack)-1]] < T[i] {
			result[stack[len(stack)-1]] = i - stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}
	return result
}
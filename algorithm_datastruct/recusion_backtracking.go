package main

// 利用递归保存上下文的能力缩小问题
//leetcode 91
func numDecodings(s string) int {
	if s[0] == '0' {
		return 0
	}
	return decode(s, len(s)-1)
}

//递归的经典模式
func decode(s string, index int) int {
	// 判断输入是否合法

	// 判断当前迭代是否结束
	if index <= 0 {
		return 1
	}

	//使用递归保存上下文，缩小问题，得到小问题的解
	cur := index
	pre := index - 1
	k, m := 0, 0
	if s[index] > '0' {
		k = decode(s, index-1)
	}
	if s[pre] == '1' || (s[pre] == '2' && s[cur] <= '6') {
		m = decode(s, index-2)
	}

	//根据小问题组合成大问题的解，整合答案
	return k + m
}

//利用递归保存上下文的能力进行遍历，在遍历的过程中进行剪树，避免不必要的遍历，实现回溯
//leetcode 39
func combinationSum(candidates []int, target int) [][]int {
	result := make([][]int, 0)
	backtracking([]int{}, candidates, 0, target, result)
	return result
}

func backtracking(solution []int, candidates []int, sum int, target int, result [][]int) {
	//不符合，结束
	if sum > target {
		return
	}
	//符合，添加情况
	if sum == target {
		copyData := make([]int, len(solution))
		copy(copyData, solution)
		result = append(result, copyData) //这里需要复制的，因为一直都是在用同一个切片的内容
		return
	}
	//添加所有可能情况,只可能添加当前和后面的数字，不能添加前面的数字，防止重合
	for index := range candidates {
		solution = append(solution, candidates[index]) //加入数字
		backtracking(solution, candidates[index:], sum+candidates[index], target, result)
		solution = solution[:len(solution)-1] //没有结束弹出
	}
}

func main() {

}

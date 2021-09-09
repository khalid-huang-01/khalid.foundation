package backtrack

import (
	"testing"
)

func backtrack(candidates, chosen []int, rsl *[][]int, curSum int , target int) {
	// 终止条件
	if curSum > target {
		return
	}
	if curSum == target {
		tmp := make([]int, len(chosen))
		copy(tmp, chosen)
		*rsl = append(*rsl, tmp)
		return
	}
	// 递归
	for i, v := range candidates {
		// 入栈
		chosen = append(chosen, v)
		// 前进
		backtrack(candidates[i:], chosen, rsl ,curSum + v, target)
		// 出栈回溯
		chosen = chosen[:len(chosen)-1]
	}
}

func combinationSum(candidates []int, target int) [][]int {
	rsl := make([][]int, 0)
	backtrack(candidates, []int{}, &rsl, 0, target)
	return rsl
}

func TestBacktrack(t *testing.T) {
	array := []int{2,3,6,7}
	t.Log(combinationSum(array, 7))
}

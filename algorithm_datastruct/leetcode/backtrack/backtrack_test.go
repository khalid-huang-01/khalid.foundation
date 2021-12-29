package backtrack

import (
	"fmt"
	"testing"
)

// leetcode 39
func backtrack(candidates, chosen []int, rsl *[][]int, curSum int , target int) {
	// 终止条件; 不符合，结束，剪枝
	if curSum > target {
		return
	}
	// 符合，添加情况，并结束
	if curSum == target {
		tmp := make([]int, len(chosen))
		copy(tmp, chosen) //这里需要复制的，因为一直都是在用同一个切片的内容
		*rsl = append(*rsl, tmp)
		return
	}
	// 递归
	//添加所有可能情况,只可能添加当前和后面的数字，不能添加前面的数字，防止重合
	for i, v := range candidates {
		// 入栈 加入数字进行尝试
		chosen = append(chosen, v)
		// 前进
		backtrack(candidates[i:], chosen, rsl ,curSum + v, target)
		// 出栈回溯 清理数组，加入index元素的情况已经通过搜试验完了，将下来要试验index+1的情况，需要把index弹出来，方便后面加入index+1【这种是不使用index的情况】，体现在上面的模型图就是第2行里面{1,2,3,4} 和{2,3,4}的关系
		chosen = chosen[:len(chosen)-1]
	}
}

func combinationSum(candidates []int, target int) [][]int {
	rsl := make([][]int, 0)
	// 这里需要传入指针引用，因为在扩容的时候会出现地址变换
	backtrack(candidates, []int{}, &rsl, 0, target)
	return rsl
}

func TestBacktrack(t *testing.T) {
	array := []int{2,3,6,7}
	t.Log(combinationSum(array, 7))
}



// leetcode 1035
// 动态规划，设dp(i, j) = nums[0:i] 和 nums[0:j]之间的最大uncrossed 连线数目
// dp(i, j) 的转态转移方程如下：
// if nums[i] == nums[j] then dp(i,j) = dp(i-1,j-1)+1
// else dp(i,j) = max(dp(i,j-1), dp(i-1,j))
// 为了方便处理，让dp(i+1,j+1) 对应nums[i]和nums[j]
func maxUncrossedLines(nums1 []int, nums2 []int) int {
	size1 := len(nums1)
	size2 := len(nums2)
	dp := make([][]int, size1+1)
	for i := 0; i <= size1; i++ {
		dp[i] = make([]int, size2+1)
	}

	for i := 1; i <= size1; i++ {
		for j := 1; j <= size2; j++ {
			if nums1[i-1] == nums2[j-1] {
				dp[i][j] = dp[i-1][j-1]+1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return dp[size1][size2]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// leetcode 1514
// eP是做跨path优化的，如果之前的path到这个节点的propobility比这个大的话，其实后续就可以不用走了
// 使用剪枝加上回溯，也会超时，最大到n = 5000
func backtrack2(visited []bool, eP []float64, curP float64, maxP *float64, graph [][]float64, cur, end, n int) {
	// 判断终止条件, 已到达终点
	if cur == end {
		if curP > *maxP {
			*maxP = curP
		}
		return
	}
	if eP[cur] > curP {
		return
	} else {
		eP[cur] = curP
	}
	// 做递归
	for i := 0; i < n; i++ {
		if graph[cur][i] == 0 || visited[i] == true {
			continue
		}
		visited[i] = true
		backtrack2(visited, eP, curP * graph[cur][i], maxP, graph, i, end, n)
		// 回溯
		visited[i] = false
	}
}

func maxProbability1(n int, edges [][]int, succProb []float64, start int, end int) float64 {
	//简历邻接矩阵，值为0表示不存在连接，否则表示连接的权重
	graph := make([][]float64, n)
	for i := 0; i < n; i++ {
		graph[i] = make([]float64, n)
	}
	for i, edge := range edges {
		graph[edge[0]][edge[1]] = succProb[i]
		graph[edge[1]][edge[0]] = succProb[i]
	}

	// 回溯
	curP, maxP := 1.0, 0.0
	visited := make([]bool, n)
	eP := make([]float64, n)
	visited[start] = true
	backtrack2(visited,eP, curP, &maxP, graph, start, end, n)
	return maxP
}



func TestMap(t *testing.T) {
	graph := make([][]int, 2)
	graph[0] = append(graph[0], 0)
	fmt.Println(graph[0])
}

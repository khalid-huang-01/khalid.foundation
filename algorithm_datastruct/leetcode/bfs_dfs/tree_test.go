package bfs_dfs

import (
	"fmt"
	"sort"
)

type TreeNode struct {
	Val int
	Left *TreeNode
	Right *TreeNode
}

// leetcode 107
func levelOrderBottom(root *TreeNode) [][]int {
	rsl_mid := make([][]int, 0)
	if root == nil {
		return rsl_mid
	}
	queue := make([]*TreeNode, 0)
	queue = append(queue, root)
	var cur *TreeNode
	for len(queue) != 0 {
		num := len(queue)
		tmp := make([]int, num)
		for i := 0; i < num; i++ {
			cur = queue[0]
			queue = queue[1:]
			tmp[i] = cur.Val
			if cur.Left != nil {
				queue = append(queue, cur.Left)
			}
			if cur.Right != nil {
				queue = append(queue, cur.Right)
			}
		}
		rsl_mid = append(rsl_mid, tmp)
	}
	rsl := make([][]int, len(rsl_mid))
	for i, j := 0, len(rsl_mid) - 1; i < len(rsl_mid); i++ {
		rsl[i] = rsl_mid[j]
		j--
	}
	return rsl
}

//leetcode 332
type Airport struct {
	name string
	isVisited bool
}

func findItinerary(tickets [][]string) []string {
	// 建立邻接表，from -> []to
	graph := make(map[string][]*Airport, 0)
	for _, ticket := range tickets {
		_, ok := graph[ticket[0]]
		if !ok {
			graph[ticket[0]] = make([]*Airport, 0)
		}
		graph[ticket[0]] = append(graph[ticket[0]], &Airport{
			name:      ticket[1],
			isVisited: false,
		})
	}
	// 做排序，后面深度优先才可以做特定查找
	for key, _ := range graph {
		sort.Slice(graph[key], func(i, j int) bool {
			return graph[key][i].name < graph[key][j].name
		})
	}

	// 做回溯查找结果
	rsl := make([]string, len(tickets)+1)
	solution := make([]string, 0)
	solution = append(solution, "JFK")
	backtracking(graph, "JFK", solution, rsl)
	return rsl
}

func backtracking(graph map[string][]*Airport, from string, solution []string, rsl []string) bool {
	// 终止条件, 条目一直就可以了
	if len(solution) == len(rsl) {
		fmt.Println("rsl", solution)
		copy(rsl, solution)
		fmt.Println(rsl)
		return true
	}
	for _, to := range graph[from] {
		if to.isVisited {
			continue
		}
		solution = append(solution, to.name)
		fmt.Println(solution)
		// 成功往下走
		to.isVisited = true
		if backtracking(graph, to.name, solution, rsl) {
			return true
		}
		// 失败回溯
		solution = solution[:len(solution)-1]
		to.isVisited = false
	}
	return false
}

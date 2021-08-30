package bfs_dfs

// leetcode 200
func numIslands(grid [][]byte) int {
	rows := len(grid)
	cols := len(grid[0])
	rsl := 0
	// 4个方向
	direction := [4][2]int{{-1,0},{1,0},{0,-1},{0,1}}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if grid[i][j] == '0' {
				continue
			}
			rsl += 1
			grid[i][j] = '0'
			// 四个方向遍历
			queue := make([][]int, 0)
			queue = append(queue, []int{i,j}) // 添加当前点
			for len(queue) != 0 {
				cur := queue[0]
				queue = queue[1:]
				for i := 0; i < 4; i++ {
					newi := cur[0] + direction[i][0]
					newj := cur[1] + direction[i][1]
					if newi >= 0 && newi < rows && newj >= 0 && newj < cols && grid[newi][newj] == '1' {
						queue = append(queue, []int{newi, newj})
						grid[newi][newj] = '0' // 标记为已访问
					}
				}
			}
		}
	}
	return rsl
}

// 查找到0的就如队列
func findOrder(numCourses int, prerequisites [][]int) []int {
	// 使用下表作为key就可以了，因为是从0编码的
	inDegreeTables := make([]int, numCourses)
	outEdgesTables := make([][]int, numCourses)
	queue := make([]int, 0)
	rsl := make([]int, 0)

	// 构建表
	for _, edge := range prerequisites {
		inDegreeTables[edge[0]]++
		outEdgesTables[edge[1]] = append(outEdgesTables[edge[1]], edge[0])
	}
	// 查找入度为0的点
	for node, indegree := range inDegreeTables {
		if indegree == 0 {
			queue = append(queue, node)
		}
	}
	// 进行遍历
	var inNode int
	for len(queue) != 0 {
		inNode = queue[0]
		queue = queue[1:]
		rsl = append(rsl, inNode)
		for _, outNode := range outEdgesTables[inNode] {
			inDegreeTables[outNode]--
			if inDegreeTables[outNode] == 0 {
				queue = append(queue, outNode)
			}
		}
	}
	// 检查有没有成环
	if len(rsl) == numCourses {
		return rsl
	}
	return []int{}
}

// leetcode 785
// 二部图
func isBipartite(graph [][]int) bool {
	colors := make([]int, len(graph)) // -1,0,1, 0 means no color
	for i := 0; i < len(graph); i++ {
		if colors[i] != 0 {
			continue
		}
		// bfs
		if !isValid(graph , colors , i) {
			return false
		}
	}
	return true
}

func isValid(graph [][]int, colors []int, start int) bool {
	colors[start] = -1
	queue := make([]int, 0)
	queue = append(queue, start)
	var node int
	for len(queue) != 0 {
		node = queue[0]
		queue = queue[1:]
		for _, next := range graph[node] {
			if colors[next] == colors[node]{
				return false
			}
			if colors[next] == 0 {
				colors[next] = -colors[node]
				queue = append(queue, next)
			}
		}
	}
	return  true
}
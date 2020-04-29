// bfs
// leetcode 133
type Node struct {
	Val        int
	Neightbors []*Node
}

func cloneGraph(node *Node) *Node {
	if node == nil {
		return node
	}
	newNode := &Node{Val: node.Val}
	// q := make([]*Node, 1)
	q := make([]*Node, 0) //必须为0.如果是1的话会报错，越界
	q = append(q, node)   //入队
	visited := make(map[int]*Node)
	visited[node.Val] = newNode //标记为已访问，防止重复进入队列中进行遍历

	for len(q) != 0 {
		tmp := q[0]
		q = q[1:]

		for _, v := range tmp.Neighbors {
			//if not visit v then set clone(tmp).Neighbors = append(clone()tmp.Neighbors, v)
			if _, ok := visited[v.Val]; !ok {
				visited[v.Val] = &Node{Val: v.Val}
				q = append(q, v)
			}
			//配置好连接关系
			visited[tmp.Val].Neighbors = append(visited[tmp.Val].Neighbors, visited[v.Val])
		}
	}
	return newNode
}

// dfs

// leetcode 785
// 二部图判断，使用bfs为基础，添加着色法
// 可以是有分划的，孤岛
func isBipartitle(graph [][]int) bool {
	// 变量初始化
	nums := len(graph)
	color := make([]int, nums) //-1 红色，0 未着色，1 蓝色
	q := make([]int, 0)        //队列
	for index := range graph {
		if color[index] != 0 {
			continue
		}
		color[index] = -1
		q = append(q, index)
		for len(q) != 0 {
			tmp := q[0]
			q = q[1:]
			for _, node := range graph[tmp] {
				if color[node] == color[tmp] {
					return false
				}
				if color[node] == 0 {
					color[node] = -color[tmp]
					q = append(q, node)
				}
			}
		}
	}
	return true
}

// dfs + 着色法
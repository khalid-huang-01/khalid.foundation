import "fmt"

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
// 可以是有分划的，孤岛, 比方说1与3相连，2与4相连，其实也是符合的
func isBipartitle(graph [][]int) bool {
	// 变量初始化
	nums := len(graph)
	color := make([]int, nums) //-1 红色，0 未着色，1 蓝色
	q := make([]int, 0)        //队列
	for index := range graph { //处理分划与孤岛的情况，也就是1只与3连接，2只与4连接
		// 表明已经属于一个区域，不用再区分了
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

//最短路径
//leetcode :743、787、1334:https://blog.csdn.net/qq_41807225/article/details/104145055
//Djjkstra https://blog.csdn.net/WYwysun/article/details/81878688: 迪杰斯特拉(Dijkstra)算法是典型最短路径算法，用于计算一个节点到其他节点的最短路径。 它的主要特点是以起始点为中心向外层层扩展(广度优先搜索思想)，直到扩展到终点为止
// leetcode 743
// 最短路径集合的最大值
// 下面的解法是基于dijkstra算法的实现，本质在于不断地去更新距离表，更新的手段其实是利用了优先队列的思想【依据全局的信息进行更新距离表】，但是是用数组实现了优先队列【遍历数据找出最值】，优先的思路是使用优先队列
func networkDelayTime(times [][]int, N int, K int) int {
	// times 转换成链接列表
	graph := make([][]int, N+1)
	cost := make([][]int, N+1)
	for _, value := range times {
		graph[value[0]] = append(graph[value[0]], value[1])
		cost[value[0]] = append(cost[value[0]], value[2])
	}

	//dijkstra算法实现
	finalDistance := make([]int, N+1)
	curDistance := make([]int, N+1)
	//init
	for i := 1; i <= N; i++ {
		finalDistance[i] = -1
		curDistance[i] = -1
	}
	finalDistance[K] = 0
	size := len(graph[K])
	for i := 0; i < size; i++ {
		curDistance[graph[K][i]] = cost[K][i]
	}
	//迭代N-1轮
	for i := 1; i < N; i++ {
		//从curDistance中挑选最小的
		var minVal int = -1
		var minIndex int = -1
		for j := 1; j <= N; j++ {
			//只从还没有获取到最小值的节点集合中查找
			if finalDistance[j] != -1 {
				continue
			}
			//需要排除掉当前无限远的节点
			if minVal == -1 || (curDistance[j] != -1 && minVal > curDistance[j]) {
				minVal = curDistance[j]
				minIndex = j
			}
		}
		fmt.Println("minIndex: ", minIndex)

		//更新finalDistance和curDistance
		finalDistance[minIndex] = minVal
		for j := 1; j <= N; j++ {
			if finalDistance[j] != -1 {
				continue
			}
			//判断有从minIndex->j的边
			dis := getDistance(graph, cost, minIndex, j)
			if dis != -1 && (curDistance[j] == -1 || curDistance[j] > minVal+dis) {
				curDistance[j] = minVal + dis
			}
		}
	}

	var result int
	//返回finalDistance的最大值; 如果存在一个无限大的值（-1），那么表示有两个节点不可达，返回-1
	for i := 1; i <= N; i++ {
		if finalDistance[i] == -1 {
			result = -1
			break
		}
		if result < finalDistance[i] {
			result = finalDistance[i]
		}
	}
	return result
}

func getDistance(graph [][]int, cost [][]int, from int, target int) int {
	size := len(graph[from])
	for i := 0; i < size; i++ {
		if graph[from][i] == target {
			return cost[from][i]
		}
	}
	return -1
}

//距离表下标为i的值j表示到当前为止从K到i的距离为j
// 使用距离表，距离表的更新依赖于与邻居的距离，用邻居的距离更新完成之后，就可以把邻居放入队列中
func networkDelayTime(times [][]int, N int, K int) int {
	graph := make([][]int, N+1)
	cost := make([][]int, N+1)
	for _, value := range times {
		graph[value[0]] = append(graph[value[0]], value[1])
		cost[value[0]] = append(cost[value[0]], value[2])
	}

	q := make([]int)
	distance := make([]int, N+1)
	for i := 1; i <= N; i++ {
		distance[i] = -1
	}
	//入队
	q = append(q, K)
	distance[K] = 0

	for len(q) != 0 {
		cur := q[0]
		q = q[1:]

		//根据当前邻居的信息+当前距离表的信息来更新距离表
		for index, neighbor := graph[cur] {
			if distance[neighbor] == -1 || distance[neighbor] > distance[cur] + cost[cur][index] {
				//distance[cur] : K -> cur 
				//cost[cur][neighbor]: cur->neighbor
				distance[neighbor] = distance[cur] + cost[cur][index]
				q = append(q, neighbor)
			}
		}
	}
	var result int
	for i := 1; i <= N; i++ {
		if result < distance[i] {
			result = distance[i]
		}
	}
	return result
}
package main

import (
	"fmt"
)

// bfs
// leetcode 133
type Node struct {
	Val       int
	Neighbors []*Node
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

// 距离表下标为i的值j表示到当前为止从K到i的距离为j
// 使用距离表，距离表的更新依赖于与邻居的距离，用邻居的距离更新完成之后，就可以把邻居放入队列中
func networkDelayTime_2(times [][]int, N int, K int) int {
	graph := make([][]int, N+1)
	cost := make([][]int, N+1)
	for _, value := range times {
		graph[value[0]] = append(graph[value[0]], value[1])
		cost[value[0]] = append(cost[value[0]], value[2])
	}

	q := make([]int, 0)
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
		for index, neighbor := range graph[cur] {
			if distance[neighbor] == -1 || distance[neighbor] > distance[cur]+cost[cur][index] {
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

// 按上面的方式构建一个：https://zhuanlan.zhihu.com/p/33162490

// 输入一个邻接矩阵，返回一个全局最短路径图，如果不存在路径，使用-1表示
func floyd(graph [][]int) [][]int {
	n := len(graph)
	// 复制出来
	dist := make([][]int, n)
	for key := range graph {
		dist[key] = make([]int, n)
		copy(dist[key], graph[key])
	}
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if dist[i][k] == -1 || dist[k][j] == -1 {
					continue
				}
				// 如果是无穷大距离/没有通路，或者有通路但是值更小就更新
				if dist[i][j] == -1 || dist[i][j] > dist[i][k]+dist[k][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}
	return dist
}

// 传入一个邻接矩阵和起始坐标，求从起始坐标到其他顶点的最短距离
func dijkstra(graph [][]int, source int) []int {
	n := len(graph)
	finished := make([]bool, n) // 用于标记是否已经完成最短距离寻找
	dist := make([]int, n)
	// 自身先标记为完成，并更新距离
	finished[source] = true
	copy(dist, graph[source]) // 更新距离

	// 进行迭代n-1轮，其实最后一轮可以不迭代的，也就是只要n-2
	for i := 1; i < n; i++ {
		// 寻找还没有找到最短距离且当前距离最短的
		minIndex := -1
		minDis := -1
		for j := 0; j < n; j++ {
			// 只从没有找到最短距离的里面找，如果当前minIndex还没有赋值，或者是出现小于的情况
			if finished[j] == false && (minIndex == -1 ||  (dist[j] != -1 && minDis > dist[j])) {
				minIndex = j
				minDis = dist[j]
			}
		}
		finished[minIndex] = true
		// 根据获取到的当前最短距离的点，作为中间点，进行更新
		for j := 0; j < n; j++ {
			// 已经找到的不再找，如果那个中间点到目标的距离本身就是无穷大，就不用再继续了
			if finished[j] == true || graph[minIndex][j] == -1{
				continue
			}
			if dist[j] == -1 || dist[j] > dist[minIndex] + graph[minIndex][j] {
				dist[j] = dist[minIndex] + graph[minIndex][j]
			}
		}
	}
	return dist
}

func main() {
	// 构造邻接矩阵
	n := 5
	graph := make([][]int, n)
	for i := 0; i < n; i++ {
		graph[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				graph[i][j] = 0
			} else {
				graph[i][j] = -1 // 怎么表示无穷大呢
			}
		}
	}
	graph[0][1] = 10
	graph[1][0] = 10
	graph[0][3] = 30
	graph[3][0] = 30
	graph[0][4] = 100
	graph[4][0] = 100
	graph[1][2] = 50
	graph[2][1] = 50
	graph[2][3] = 20
	graph[3][2] = 20
	graph[2][4] = 10
	graph[4][2] = 10
	graph[3][4] = 60
	graph[4][3] = 60
	//fmt.Println(graph)
	fmt.Println(floyd(graph))
	fmt.Println(dijkstra(graph, 0))

}

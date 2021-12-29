package dynamic

import (
	"container/heap"
	"fmt"
	"testing"
)

func numTrees(n int) int {
	dp := make([]int, n+1)
	dp[0] = 1
	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			dp[i] += dp[j] * dp[i-j-1]
		}
	}
	return dp[n]
}

// 最短路径问题

func TestShortestPath(t *testing.T) {
	graph := [][]int{
		{
			0, 10, -1, 30, 100,
		}, {
			-1, 0, 50, -1, 10,
		}, {
			-1, -1, 0, -1, 10,
		}, {
			-1, -1, 20, 0, -1,
		}, {
			-1, -1, -1, 60, 0,
		},
	}
	//t.Log(floyd(graph))
	t.Log(dijkstra(graph, 0))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 多源最短路径 d(i,j) = min{d(i,j), d(i,u)+d(u,j)}
func floyd(graph [][]int) [][]int {
	n := len(graph)
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, len(graph[i]))
		copy(dist[i], graph[i])
	}
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if dist[i][k] == -1 || dist[k][j] == -1 {
					continue
				}
				if dist[i][j] == -1 || dist[i][j] > dist[i][k]+dist[k][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}
	return dist
}

// 单源最短路径, 公式要比上面的少一个变量
// d(i) = min{d(i), d(j)+u(j,i)} //
// 基于点来遍历实现上面的公式
func dijkstra(graph [][]int, source int) []int {
	n := len(graph)
	dist := make([]int, n)
	copy(dist, graph[source])
	finished := make([]bool, n)
	finished[source] = true
	for i := 1; i < n; i++ {
		// 找到有最小距离的点
		minIndex, minValue := -1, -1 // 最大值
		for j := 0; j < n; j++ {
			if finished[j] || dist[j] == -1 { // 如果已经找到最小值的点，或者是无穷大的
				continue
			}
			if minValue == -1 || minValue > dist[j] {
				minValue = dist[j]
				minIndex = j
			}
		}

		// 做松弛
		// d(i) = min{d(i), d(j)+u(j,i)} //
		finished[minIndex] = true
		for j := 0; j < n; j++ {
			if finished[j] || graph[minIndex][j] == -1 {
				continue
			}
			if dist[j] == -1 || dist[j] > dist[minIndex]+graph[minIndex][j] {
				dist[j] = dist[minIndex] + graph[minIndex][j]
			}
		}
	}
	return dist
}

// leetcode 133
type Node struct {
	Val       int
	Neighbors []*Node
}

func cloneGraph(node *Node) *Node {
	if node == nil {
		return node
	}
	isVisited := make(map[*Node]struct{})
	queue := make([]*Node, 0)
	// 广度优先搜索
	queue = append(queue, node)
	isVisited[node] = struct{}{}
	start := &Node{
		Val: node.Val,
	}
	queue1 := make([]*Node, 0)
	queue1 = append(queue1, start)
	var cur, newCur *Node
	for len(queue) != 0 {
		cur = queue[0]
		queue = queue[1:]
		newCur = queue1[0]
		queue1 = queue1[1:]
		for _, neighbor := range cur.Neighbors {
			newNeighbor := &Node{
				Val: neighbor.Val,
			}
			newCur.Neighbors = append(newCur.Neighbors, newNeighbor)
			if _, ok := isVisited[neighbor]; !ok {
				isVisited[neighbor] = struct{}{}
				queue = append(queue, neighbor)
				queue1 = append(queue1, newNeighbor)
			}
		}
	}
	return start
}

func Test_Clone_Graph(t *testing.T) {
	var n1 = Node{Val: 1}
	var n2 = Node{Val: 2}
	var n3 = Node{Val: 3}
	var n4 = Node{Val: 4}

	n1.Neighbors = []*Node{&n2, &n4}
	n2.Neighbors = []*Node{&n1, &n3}
	n3.Neighbors = []*Node{&n2, &n4}
	n4.Neighbors = []*Node{&n1, &n3}

	r := cloneGraph(&n1)
	PrintNode(r)
	PrintNode(r.Neighbors[0])
	PrintNode(r.Neighbors[0].Neighbors[1])
	PrintNode(r.Neighbors[1])
}

func PrintNode(node *Node) {
	fmt.Println("val: ", node.Val)
	fmt.Println("neighbors: ", len(node.Neighbors))
	for _, neighbor := range node.Neighbors {
		fmt.Println(neighbor.Val, " ")
	}
}

// leetcode 1031:
// 这个题目就是典型的控制变量法，把其中一个数组固定，另外一个动态找最大的
func maxSumTwoNoOverlap(nums []int, firstLen int, secondLen int) int {
	size := len(nums)
	dpf := make([]int, size+1)
	dps := make([]int, size+1)
	dp := make([]int, size+1)

	for i := firstLen - 1; i < size; i++ {
		// i是最后一个位置
		dpf[i+1] = max(accumulate(nums, i-firstLen+1, firstLen), dpf[i])
	}
	for i := secondLen - 1; i < size; i++ {
		dps[i+1] = max(accumulate(nums, i-secondLen+1, secondLen), dps[i])
	}

	// dp的长度最少也要是firstLen+secondLen
	for i := firstLen + secondLen - 1; i < size; i++ {
		// 在0-i之间查早最大的不重叠和：
		//		固定firstLen，动态使用最大的second
		sum1 := accumulate(nums, i-firstLen+1, firstLen) + dps[i-firstLen+1]
		//		固定secondLen, 动态使用最大的first
		sum2 := accumulate(nums, i-secondLen+1, secondLen) + dpf[i-secondLen+1]
		dp[i+1] = max(sum1, sum2)
		dp[i+1] = max(dp[i], dp[i+1])
	}
	return dp[size]
}

func accumulate(nums []int, start, len int) int {
	sum := 0
	for i := 0; i < len; i++ {
		sum += nums[start+i]
	}
	return sum
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// leetcode 1514
// 动态规划，也就是 Dijktra 算法
// dijkstra + 邻接表 （不可用矩阵）
// 算法没问题，在n = 10000 超时了
func maxProbability1(n int, edges [][]int, succProb []float64, start int, end int) float64 {
	// 构建邻接矩阵，就是一个n*n的矩阵，值可以表示权重，用矩阵太大，超过内存，
	// 使用稀疏的表达方法，也就是邻接表的方式每个节点，接一个列表，列表里面时自己的连接端点的结合
	graph := make([]map[int]float64, n)
	for i := 0; i < n; i++ {
		graph[i] = make(map[int]float64)
	}
	for i, edge := range edges {
		graph[edge[0]][edge[1]] = succProb[i]
		graph[edge[1]][edge[0]] = succProb[i]
	}

	// 初始dijkstra参数
	singleProb := make([]float64, n)
	finished := make([]bool, n)
	for vertex, weight := range graph[start] {
		singleProb[vertex] = weight
	}
	finished[start] = true

	// 迭代更新单源距离表
	// 动规方程：p[i] = max(p[i], p[j]*graph(j, i)) // 这个表示时一个迭代了
	//      p[i]表示start到i的最大概率
	//		start节点到i的最大概率是：等于自己本身当前的值与start到j的概率与j到i的概率积；就是先到start到j，然后从j到i
	var maxIndex int
	var maxPro float64
	for i := 0; i < n ; i++ {
		// 查找当前最大pro的endpoint
		maxPro = 0
		maxIndex = -1
		// 这一块需要修改为优先队列，才能过估计
		for j := 0; j < n; j++ {
			if !finished[j] && (maxPro == 0 || maxPro < singleProb[j]) {
				maxPro = singleProb[j]
				maxIndex = j
			}
		}

		// fmt.Println(maxIndex)
		if maxIndex == -1 || maxIndex == end {
			break
		}
		finished[maxIndex] = true

		// 基于这个endpoint来松弛其他的endpoint
		for j := 0; j < n ; j++ {
			if !finished[j] && singleProb[j] < singleProb[maxIndex] * graph[maxIndex][j] {
				singleProb[j] = singleProb[maxIndex] * graph[maxIndex][j]
			}
		}

		// fmt.Println(singleProb[end])
	}

	// 返回结果
	return singleProb[end]
}

// https://blog.csdn.net/Ratina/article/details/86774613
// leetcode 1514
// Dijkstra + 优先队列优化
// 动规方程：p[i] = max(p[i], p[j]*graph(j, i)) // 这个表示时一个迭代了
//      p[i]表示start到i的最大概率
//		start节点到i的最大概率是：等于自己本身当前的值与start到j的概率与j到i的概率积；就是先到start到j，然后从j到i
func maxProbability(n int, edges [][]int, succProb []float64, start int, end int) float64 {
	// 稀疏，构建邻接表
	graph := make([]map[int]float64, n)// 第一个map是节点，每个节点对于一个map，map是一个连接节点和概率的kv
	for i := 0; i < n; i++ {
		graph[i] = make(map[int]float64, 0)
	}

	for i, edge := range edges {
		graph[edge[0]][edge[1]] = succProb[i]
		graph[edge[1]][edge[0]] = succProb[i]
	}

	// 初始化dijkstra参数
	h := &Infos{
		Info{
			node: start,
			pro:  1,
		},
	}
	heap.Init(h)
	finished := make([]bool, n)

	// 执行迭代
	for h.Len() > 0 {
		cur := heap.Pop(h).(Info)

		if cur.node == end {
			return cur.pro
		}

		// 这里不断迭代，会把相同的节点里面，不是大得过滤掉
		if finished[cur.node] {
			continue
		}
		finished[cur.node] = true

		//p[i] = max(p[i], p[j]*graph(j, i)), 这里不用比较都放入到优先队列里面，通过317行过滤就可以了
		// 针对cur的边做松弛
		for neighbor, pro := range graph[cur.node] {
			if finished[neighbor] {
				continue
			}
			heap.Push(h, Info{
				node: neighbor,
				pro: cur.pro * pro,
			})
		}
	}
	return 0
}

type Info struct {
	node int
	pro float64
}

type Infos []Info
func (p *Infos) Push(x interface{}) {
	*p = append(*p, x.(Info))
}


func (p *Infos) Pop() interface{} {
	old := *p
	popped := old[len(old)-1]
	*p = old[:len(old)-1]

	return popped
}

func (p Infos) Len() int {
	return len(p)
}
func (p Infos) Less(i, j int) bool {
	return p[i].pro > p[j].pro
}

// Swap swaps the elements with indexes i and j.
func (p Infos) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func TestM(t *testing.T)  {
	p := Infos{}
	p = append(p, Info{
		node: 2,
		pro:  0.2,
	})
	p = append(p, Info{
		node: 1,
		pro:  0.1,
	})
	p = append(p, Info{
		node: 4,
		pro:  0.3,
	})
	p = append(p, Info{
		node: 5,
		pro:  0.1,
	})

	heap.Init(&p)
	temp := heap.Pop(&p).(Info)
	fmt.Println(temp)

}


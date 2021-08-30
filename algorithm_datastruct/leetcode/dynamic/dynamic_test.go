package dynamic

import (
	"fmt"
	"testing"
)

func numTrees(n int) int {
	dp := make([]int ,n+1)
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
	graph := [][]int {
		{
			0,10,-1,30,100,
		}, {
			-1,0,50,-1,10,
		}, {
			-1,-1,0,-1,10,
		}, {
			-1,-1,20,0,-1,
		}, {
			-1,-1,-1,60,0,
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
				if dist[i][j] == -1 || dist[i][j] > dist[i][k] + dist[k][j] {
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
			if finished[j] || dist[j] == -1{ // 如果已经找到最小值的点，或者是无穷大的
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
			if dist[j] == -1 || dist[j] > dist[minIndex] + graph[minIndex][j] {
				dist[j] = dist[minIndex] + graph[minIndex][j]
			}
		}
	}
	return dist
}

// leetcode 133
type Node struct {
	Val int
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
		Val:       node.Val,
	}
	queue1 := make([]*Node ,0)
	queue1 = append(queue1, start)
	var cur, newCur *Node
	for len(queue) != 0 {
		cur = queue[0]
		queue = queue[1:]
		newCur = queue1[0]
		queue1 = queue1[1:]
		for _, neighbor := range cur.Neighbors {
			newNeighbor := &Node{
				Val:       neighbor.Val,
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
	var n1 = Node{ Val : 1}
	var n2 = Node{ Val : 2}
	var n3 = Node{ Val : 3}
	var n4 = Node{ Val : 4}

	n1.Neighbors = []*Node{ &n2, &n4 }
	n2.Neighbors = []*Node{ &n1, &n3 }
	n3.Neighbors = []*Node{ &n2, &n4 }
	n4.Neighbors = []*Node{ &n1, &n3 }

	r := cloneGraph(&n1)
	PrintNode(r)
	PrintNode(r.Neighbors[0])
	PrintNode(r.Neighbors[0].Neighbors[1])
	PrintNode(r.Neighbors[1])
}

func PrintNode(node *Node) {
	fmt.Println("val: ", node.Val)
	fmt.Println("neighbors: ",len(node.Neighbors))
	for _, neighbor := range node.Neighbors {
		fmt.Println(neighbor.Val, " ")
	}
}
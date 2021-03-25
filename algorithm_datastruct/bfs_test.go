package main

import (
	"testing"
)

// -------------- leetcode 127
func ladderLength(beginWord string, endWord string, wordList []string) int {
	// 生成邻接表
	wordList = append(wordList, beginWord)
	graph := InitGraph(wordList)
	// 根据邻接表，进行广度优先搜索，每层结束添加一个空串
	queue := make([]string, 0)
	visited := make(map[string]bool, 0)
	queue = append(queue, beginWord)
	visited[beginWord] = true
	rsl, curLayerSize := 1, 0
	var word string
	for len(queue) != 0 {
		// 进行当前层的数量，是使用层来遍历的，所以每次初始的个数就是当前层的节点数
		curLayerSize = len(queue)
		for i := 0; i < curLayerSize; i++{
			word = queue[0]
			queue = queue[1:]
			if word == endWord {
				return rsl
			}
			for _, next := range graph[word] {
				if visited[next] == false {
					queue = append(queue, next)
					visited[next] = true
				}
			}
		}
		rsl += 1
	}
	return 0
}

// 生成邻接表
func InitGraph(worldList []string) map[string][]string {
	graph := make(map[string][]string)
	for i := 0; i < len(worldList); i++ {
		for j := i + 1; j < len(worldList); j++ {
			if isChangeAble(worldList[i], worldList[j]) {
				graph[worldList[i]] = append(graph[worldList[i]], worldList[j])
				graph[worldList[j]] = append(graph[worldList[j]], worldList[i])
			}
		}
	}
	return graph
}

func isChangeAble(word1 string, word2 string) bool {
	diff := 0
	for i := 0; i < len(word1); i++ {
		if word1[i] != word2[i] {
			diff++
		}
	}
	return diff == 1
}

func TestLadderLength(t *testing.T) {
	beginWord := "hit"
	endWord := "cog"
	wordList := []string{"hot","dot","dog","lot","log","cog"}
	t.Log(ladderLength(beginWord, endWord, wordList))
}
// --------------leetcode 127

//-------------leetcode 139
func wordBreak(s string, wordDict []string) bool {
	queue := make([]int, 0)
	wordMap := make(map[string]bool, 0)
	visited := make(map[int]bool, 0)
	for _, word := range wordDict {
		wordMap[word] = true
	}
	queue = append(queue, 0)
	visited[0] = true
	var start int
	var result bool
	var word string
	size := len(s)
	for len(queue) != 0 {
		start = queue[0]
		queue = queue[1:]
		for end := start + 1; end <= size; end++ {
			if visited[end] {
				continue
			}
			word = s[start:end]
			if wordMap[word] {
				if end == size {
					result = true
					break
				}
				queue = append(queue, end)
				visited[end] = true
			}
		}
	}
	return result
}

func TestWordBreak(t *testing.T) {
	s := "catsandog"
	wordDict := []string{"cats", "dog", "sand", "and", "cat"}
	print(wordBreak(s, wordDict))
}


func solve(board [][]byte)  {
	if len(board) == 0 {
		return
	}
	// 从边缘获取到全部为O的坐标放入队列中，并标记为已访问
	rows := len(board)
	cols := len(board[0])
	queue := make([][]int, 0) // 二维数组，第二维用于表示节点的坐标信息
	visited := make([][]bool, rows)
	for i := 0; i < rows; i++ {
		visited[i] = make([]bool, cols)
	}

	for i := 0; i < cols; i++ {
		//第一行
		if board[0][i] == 'O' {
			queue = append(queue, []int{0, i})
			visited[0][i] = true
		}
		// 最后一行
		if board[rows-1][i] == 'O' {
			queue = append(queue, []int{rows-1, i})
			visited[rows-1][i] = true
		}
	}
	for i := 1; i < rows; i++ {
		//第一列
		if board[i][0] == 'O' {
			queue = append(queue, []int{i, 0})
			visited[i][0] = true
		}
		if board[i][cols-1] == 'O' {
			queue = append(queue, []int{i, cols-1})
			visited[i][cols-1] = true
		}
	}
	// bfs进行处理，标记所有相连的O为已访问

	dir := [4][2]int { {-1,0},{1,0},{0,1},{0,-1} } // 四个方向
	for len(queue) != 0 {
		node := queue[0]
		queue = queue[1:]
		//进四个方向进行查看
		x, y := node[0], node[1]
		var nx, ny int
		for i := 0; i < 4; i++ {
			nx = x + dir[i][0]
			ny = y + dir[i][1]
			if nx >= 0 && nx < rows && ny >= 0 && ny < cols && board[nx][ny] == 'O' && !visited[nx][ny] {
				queue = append(queue, []int{nx, ny})
				visited[nx][ny] = true
			}
		}
	}
	// 将为O且未访问的点替换为X
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if board[i][j] == 'O' && !visited[i][j] {
				board[i][j] = 'X'
			}
		}
	}
}

type Point struct {
	x, y int
}
// leetcode 200
func numIslands1(grid [][]byte) int {
	result := 0
	rows := len(grid)
	cols := len(grid[0])
	// 尽量把1变成0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if grid[i][j] == '1' {
				result += 1
				grid[i][j] = '0'
				queue := make([]Point, 0)
				queue = append(queue, Point{x:i,y:j})
				for len(queue) != 0 {
					point := queue[0]
					queue = queue[1:]
					if point.x - 1 >= 0&& grid[point.x-1][point.y] == '1' {
						queue = append(queue, Point{x:point.x-1,y:point.y})
						grid[point.x-1][point.y]='0'
					}
					if point.x + 1 < rows&& grid[point.x+1][point.y] == '1' {
						queue = append(queue, Point{x:point.x+1,y:point.y})
						grid[point.x+1][point.y]='0'
					}
					if point.y - 1 >= 0&& grid[point.x][point.y-1] == '1' {
						queue = append(queue, Point{x:point.x,y:point.y-1})
						grid[point.x][point.y-1]='0'
					}
					if point.y + 1 < cols&& grid[point.x][point.y+1] == '1' {
						queue = append(queue, Point{x:point.x,y:point.y+1})
						grid[point.x][point.y+1]='0'
					}
				}
			}
		}
	}
	return result
}
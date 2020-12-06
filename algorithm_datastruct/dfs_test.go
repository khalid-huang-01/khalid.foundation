package main

import (
	"fmt"
	"sort"
	"testing"
)

//import (
//	"sort"
//)
// leetcode 332
type Airport struct {
	airport string
	visited bool
}

type Airports []*Airport

func (as Airports) Len() int {
	return len(as) 
}

func (as Airports) Swap(i, j int) {
	as[i], as[j] = as[j], as[i]
}

func (as Airports) Less(i, j int) bool {
	return as[i].airport < as[j].airport
}


var ticketsCount int

func findItinerary(tickets [][]string) []string {
	// 组织起邻接表
	ticketsCount = len(tickets)

	graph := make(map[string]Airports, 0)
	for _, ticket := range tickets {
		_, ok := graph[ticket[0]]
		if !ok {
			graph[ticket[0]] = make(Airports, 0)
		}
		graph[ticket[0]] = append(graph[ticket[0]], &Airport{visited: false, airport: ticket[1]})
	}

	// 对每个都表都进行下排序
	for key,_ := range graph {
		sort.Sort(graph[key])
	}
	printGraph(graph)

	cur := "JFK"
	// 使用回溯法进行深度优先搜索
	rsl := make([]string, ticketsCount+1)
	solution := make([]string, 0)
	solution = append(solution, cur)
	backtracking1(graph, cur, solution, rsl)
	return rsl
}

func backtracking1(graph map[string]Airports, cur string, solution []string, rsl []string) bool {
	// 结束条件判断
	// 找到一个就可以了
	if len(solution) == ticketsCount + 1{
		copy(rsl, solution)
		fmt.Println("solution: ", solution)
		return true
	}
	// 进行回溯
	for _, to := range graph[cur] {
		if to.visited {
			continue
		}
		solution = append(solution, to.airport)
		to.visited = true
		if backtracking1(graph, to.airport, solution, rsl) {
			return true
		}
		solution = solution[:len(solution)-1]
		to.visited = false

	}
	return false

}

func printGraph(graph map[string]Airports) {
	fmt.Println("----------------")
	for key := range graph {
		for _, airport := range graph[key] {
			fmt.Print(airport.airport, " ")
		}
		fmt.Println("")
	}
	fmt.Println("----------------")
}

func TestFindItinerary(t *testing.T) {
	tickets := [][]string {
		{"JFK","SFO"},
		{"JFK","ATL"},
		{"SFO","ATL"},
		{"ATL","JFK"},
		{"ATL","SFO"},
	}
	fmt.Println(findItinerary(tickets))
}

func testslice(s []string) {
	s = s[:len(s)-1]
}

func TestSlice(t *testing.T) {
	s := []string {"1", "2"}
	testslice(s)
	fmt.Println(s)
}
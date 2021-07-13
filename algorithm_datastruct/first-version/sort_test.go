package main

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

// 每次都把最大的数据冒到坐标最大的位置
func bubbleSort(data []int) {
	size := len(data)
	for i := 0; i < size; i++ { //迭代的次数
		for j := 0; j < size-i-1; j++ { //当前迭代需要比较的次数
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
}

func insertSort(data []int) {
	size := len(data)
	for i := 1; i < size; i++ {
		//在排好序的数组中查找需要这个新数字的位置，每次比较，成功就往后移动
		tmp := data[i] //把i位置的data位置空出来，方便后面存放数据
		var j int
		for j = i - 1; j >= 0; j-- {
			if tmp < data[j] {
				data[j+1] = data[j]
			} else {
				break
			}
		}

		data[j+1] = tmp
	}
}

func _mergeSort(data []int, low int, high int) {
	//合法性检查和终止条件
	if low >= high {
		return
	}
	//分解，归并小解答案
	middle := low + (high-low)/2
	_mergeSort(data, low, middle)
	_mergeSort(data, middle+1, high)
	merge(data, low, middle, high)
}

//赋值的手法高明
func merge1(data []int, low int, middle int, high int) {
	tmp := make([]int, len(data))
	copy(tmp, data)
	k, i, j := low, low, middle+1
	//k代表需要为当data赋值的下标
	for k <= high {
		if i > middle {
			data[k] = tmp[j]
			j++
		} else if j > high {
			data[k] = tmp[i]
			i++
		} else if tmp[i] > tmp[j] {
			data[k] = tmp[j]
			j++
		} else {
			data[k] = tmp[i]
			i++
		}
		k++
	}
}

// 使用哨兵简化编程，主要是为了避免特殊情况的处理，在两个子串最后都加一个无穷大的数字，这样就可以避免有子串遍历完的情况
const INT_MAX = int(^uint(0) >> 1) // 全部取反再右移一位，就只有第一位是0了
const INT_MIN = ^INT_MAX
func merge(data []int, low int, middle int, high int) {
	tmpLeft := make([]int, middle-low+2) // middle-low+1+ 1 //最后都补1
	tmpRight := make([]int, high-middle+1) // high - (middle+1) + 1 + 1
	copy(tmpLeft, data[low:middle+1])
	tmpLeft[len(tmpLeft)-1] = INT_MAX
	copy(tmpRight, data[middle+1:high+1])
	tmpRight[len(tmpRight)-1] = INT_MAX

	for k,i,j := low, 0, 0; k <= high; k++ {
		if tmpLeft[i] < tmpRight[j] {
			data[k] = tmpLeft[i]
			i++
		} else {
			data[k] = tmpRight[j]
			j++
		}
	}
}

func mergeSort(data []int) {
	low := 0
	high := len(data) - 1
	_mergeSort(data, low, high)
}

func TestMergeSort(t *testing.T) {
	data := []int{3,1,5,2,7,11,4}
	mergeSort(data)
	t.Log(data)
}

func _quickSort(data []int, low int, high int) {
	//终止条件
	if low >= high {
		return
	}
	//划分
	index := partition(data, low, high)

	_quickSort(data, low, index-1)
	_quickSort(data, index+1, high)

}

func partition(data []int, low int, high int) int {
	//随机选择一个数做为基准值，并与high位置交互 => high作为一个基准值了
	swap(data, low+rand.Intn(high-low+1), high)

	//换位,双指针
	var i, j int
	// 思想其实就是把列表划分为两个区，low到i-1区间表示已选择的比data[high]小的数字，i到j-1之间的是比data[high]大的数字，j以上的是没有判断的
	// 算法的目的就是在没有判断的里面找到比data[high]小的，然后与data[i]交互，放入小区数字，再把i往前移一位，到大区数字，恢复秩序
	// i是第一个比基准值大的数
	// j是用于判断的数
	for i, j = low, low; j < high; j++ {
		if data[j] < data[high] {
			swap(data, i, j)
			i++
		}
	}

	swap(data, i, high)
	return i

}

func swap(data []int, i int, j int) {
	data[i], data[j] = data[j], data[i]
}

func quickSort(data []int) {
	low := 0
	high := len(data) - 1
	_quickSort(data, low, high)
}

// 堆排序：https://www.cnblogs.com/xingyunshizhe/p/11311754.html

// 使用内置的排序算法 https://www.jianshu.com/p/4bf3c94a15a6
type person struct {
	Name string
	Age  int
}

func TestSort(t *testing.T) {
	// data := []int{4, 1, 2, 6, 8, 3}
	// bubbleSort(data)
	// insertSort(data)
	// mergeSort(data)
	// quickSort(data)
	// fmt.Println(data)
	a := []person{
		{
			Name: "AAA",
			Age:  55,
		},
		{
			Name: "BBB",
			Age:  22,
		},
		{
			Name: "CCC",
			Age:  0,
		},
		{
			Name: "DDD",
			Age:  22,
		},
		{
			Name: "EEE",
			Age:  11,
		},
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i].Age < a[j].Age
	})
	fmt.Println(a)
}

// 拓扑排序 210
// Topological Order，拓扑排序
func findOrder(numCourses int, prerequisites [][]int) []int {
	// 入度表和出边表
	inDegree := make([]int, numCourses)
	outEdges := make([][]int, numCourses)

	// 根据prerequisites进行根据入度表和出边表
	for _, v := range prerequisites {
		inDegree[v[0]]++
		outEdges[v[1]] = append(outEdges[v[1]], v[0])
	}

	queue := make([]int, 0)
	rsl := make([]int , 0)
	// 获取入度为0的点，加入队列
	for i, v := range inDegree {
		if v == 0 {
			queue = append(queue, i)
		}
	}

	var node int
	for len(queue) != 0 {
		node = queue[0]
		queue = queue[1:] // 将队头出队
		// 消除以队头为源点的边
		for _, v := range outEdges[node] {
			inDegree[v]--
			if inDegree[v] == 0 {
				queue = append(queue, v)
			}
		}
		// 添加结果
		rsl = append(rsl, node)
	}
	// 检查有没有成环，如果成环，直接返回空
	if len(rsl) != numCourses {
		rsl = []int{}
	}
	return rsl
}

type Group struct {
	items map[int]int // 建立原始数据到表中的下标映射关系
	items2 map[int]int //建立上面的反向关系
	count int //计数
	inDegree []int
	outEdges [][]int
}

func (g *Group) Add(num int) {
	g.items[num] = g.count
	g.count++
}

// 拓扑排序 1203
// nums节点数量， inDegree每个节点的入度，outEdges每个节点的出边
func topologicalSort(nums int, inDegree []int, outEdges [][]int) []int {
	queue := make([]int, 0)
	for i, v := range inDegree {
		if v == 0 {
			queue = append(queue, i)
		}
	}
	var node int
	rsl := make([]int,0)
	for len(queue) != 0 {
		node = queue[0]
		rsl = append(rsl, node)
		queue = queue[1:]
		// 删除以这个为出发点的边
		for _, v := range outEdges[node] {
			inDegree[v]--
			if inDegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
	if len(rsl) != nums {
		rsl = []int{}
	}
	return rsl
}

func sortItems(n int, m int, group []int, beforeItems [][]int) []int {
	// 使用 group建立组别，为-1的，都给新的组号
	groups := make(map[int]*Group)
	maxGroupIndex := m - 1
	for i, v := range group {
		if v == -1 {
			maxGroupIndex++
			groups[maxGroupIndex] = &Group{
				items: map[int]int{},
			}
			groups[maxGroupIndex].Add(i)
			group[i] = maxGroupIndex // 建立反向查找
		} else {
			_, ok := groups[v]
			// 没有初始化过的话，要进行初始化
			if !ok {
				groups[v] = &Group{
					items: map[int]int{},
				}
			}
			groups[v].Add(i)
		}
	}
	// 根据上述情况对inDegress等进行初始化
	for _, g := range groups {
		if len(g.items) > 1 {
			g.inDegree = make([]int, len(g.items))
			g.outEdges = make([][]int, len(g.items))
		}
	}


	// 使用beforeItems建立group与group之间，和group内的拓扑连接
	groupInDegree := make([]int, maxGroupIndex)
	groupOutEdges := make([][]int, maxGroupIndex)
	// 2.1 使用beforeItem进行建立组之间和组与组之间的拓扑关系
	var endGroup, startGroup int
	var endNode, startNode int
	for end, items := range beforeItems {
		for _, start := range items {
			endGroup = group[end]
			startGroup = group[start]
			fmt.Println(start, " ", startGroup)
			fmt.Println(end, " ", endGroup)

			// 建立组之间的关系
			if endGroup != startGroup {
				groupInDegree[endGroup]++
				groupOutEdges[startGroup] = append(groupOutEdges[startGroup], endGroup)
			} else {
				// 建立节点之间的关系
				endNode = groups[startGroup].items[end]
				startNode = groups[startGroup].items[start]
				groups[startGroup].inDegree[endNode]++
				groups[startGroup].outEdges[startNode] = append(groups[startGroup].outEdges[startNode], endNode)

			}
		}
	}
	// 2.2 对组与组之间进行拓扑排序
	groupTopological := topologicalSort(len(groups), groupInDegree, groupOutEdges)
	// 2.2 针对每组再进行拓扑排序
	rsl := make([]int ,0)
	for _, v := range groupTopological {
		if len(groups[v].items) == 1 {
			rsl = append(rsl, groups[v].items[0])
		} else {
			tmp := topologicalSort(len(groups[v].items), groups[v].inDegree, groups[v].outEdges)
			if len(tmp) == 0 {
				rsl = []int{}
				break
			} else {
				rsl = append(rsl, tmp...)
			}
		}
	}
	return  rsl
}


func TestSortItems(t *testing.T) {
	n := 8
	m := 2
	group := []int{-1,-1,1,0,0,1,0,-1}
	beforeItems := [][]int{
		[]int{},
		[]int{6},
		[]int{5},
		[]int{6},
		[]int{3,6},
		[]int{},
		[]int{},
		[]int{},
	}
	t.Log(sortItems(n,m,group,beforeItems))
}
func TestRange(t *testing.T) {
	arr := make([][]int, 2)
	for _, v := range arr[0] {
		fmt.Println(v)
	}
}

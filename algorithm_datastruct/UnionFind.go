package main

import "sort"

// ------------------- leetcode 684
type UnionFindSet684 struct {
	parent []int
}

func (u *UnionFindSet684) Init(nums int) {
	u.parent = make([]int, nums)
	for i := range u.parent {
		u.parent[i] = i // 初始化为指向自己
	}
}

func (u *UnionFindSet684) Find(node int) int {
	root := node
	// 找到最高级的节点
	for u.parent[root] != root {
		root = u.parent[root]
	}
	// 将其他的后辈节点都指向这个最高级节点
	son := node
	var tmp int
	for u.parent[son] != root {
		tmp = u.parent[son]
		u.parent[son] = root
		son = tmp
	}
	return root
}

// 返回false表示两个节点已经在同一个集合里面了，不需要再进行联合
func (u *UnionFindSet684) Union(node1 int, node2 int) bool {
	root1 := u.Find(node1)
	root2 := u.Find(node2)
	if root1 == root2 {
		return false
	}
	if root1 < root2 {
		u.parent[root1] = root2
	} else {
		u.parent[root2] = root1
	}
	return true
}

func findRedundantConnection(edges [][]int) []int {
	u := &UnionFindSet684{}
	nums := len(edges)
	u.Init(nums+1)
	// 遍历并获取结果
	var result []int
	for _, edge := range edges {
		if !u.Union(edge[0], edge[1]) {
			result = edge
		}
	}
	return result
	//
}


// ------------------- leetcode 684

// ------------------- leetcode 924

type UnionFindSet struct {
	Parents []int //记录每个节点的上级
	Size []int // 记录以当前结点为顶点的连通分量里面的节点有多少个（只有自己的话，值为1）
}

func (u *UnionFindSet) Init(size int) {
	u.Parents = make([]int, size)
	u.Size = make([]int, size)
	for i := 0; i < size; i++ {
		u.Parents[i] = i
		u.Size[i] = 1
	}
}

func (u *UnionFindSet) Find(node int) int {
	if u.Parents[node] == node {
		return node
	}
	root := u.Find(u.Parents[node])
	u.Parents[node] = root
	return root
}

func (u *UnionFindSet) Union(node1 int, node2 int) {
	root1 := u.Find(node1)
	root2 := u.Find(node2)
	if root1 == root2 {
		return
	}
	if root1 < root2 {
		u.Parents[root1] = root2
		u.Size[root2] += u.Size[root1] // 以root2为最顶层结点的连通分量的个数要叠加上root1的
	}
}

func minMalwareSpread(graph [][]int, initial []int) int {
	// 初始化并查集
	u := &UnionFindSet{}
	u.Init(len(graph))

	// 根据graph进行连通，生成连通分量，并记录连通分量的大小
	for i := 0; i < len(graph); i++ {
		for j := 0; j < len(graph[i]); j++ {
			if graph[i][j] == 1 {
				u.Union(i,j)
			}
		}
	}

	// 查找目标，统计每个感染节点的颜色情况
	// 先对Init进行排序
	sort.Ints(initial)
	count := make(map[int]int,0)
	for i := 0; i < len(initial); i++ {
		count[u.Find(initial[i])]++
	}
	// 1. 如果有唯一颜色的，就选择其中连通分量最大的
	ans := -1
	ansSize := -1
	root := 0
	for i := 0; i < len(initial); i++ {
		// 是唯一颜色的
		root = u.Find(initial[i])
		if count[root] == 1 && (ans == -1 || ansSize < u.Size[root]) {
			ans = initial[i]
			ansSize = u.Size[root]
		}
	}

	// 2. 如果没有唯一颜色的节点，就选择下标最小的
	if ans == -1 {
		ans = initial[0]
	}
	return ans
}
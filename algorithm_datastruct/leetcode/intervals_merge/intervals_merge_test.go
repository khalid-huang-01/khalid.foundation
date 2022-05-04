package intervals_merge

import "sort"

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// leetcode 56
func merge(intervals [][]int) [][]int {
	// 按降序排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	ans := make([][]int, 0)
	for _, interval := range intervals {
		// 如果为空，或者最后一个的右侧比当前的左侧还小，就添加新的
		if len(ans) == 0 || ans[len(ans)-1][1] < interval[1] {
			ans = append(ans, []int{interval[0], interval[1]})
		} else {
			// 更新
			ans[len(ans)-1][1] = max(ans[len(ans)-1][1], interval[1])
		}
	}
	return ans
}

// leetcode 57
func insert(intervals [][]int, newInterval []int) [][]int {
	ans := make([][]int, 0)
	size := len(intervals)
	i := 0
	// 左边不重叠的直接入
	for i < size && intervals[i][1] < newInterval[0] {
		ans = append(ans ,intervals[i])
		i++
	}
	// 中间重叠的合并
	for i < size && intervals[i][0] <= newInterval[1] {
		newInterval[0] = min(intervals[i][0], newInterval[0])
		newInterval[1] = max(intervals[i][1], newInterval[1])
		i++
	}
	ans = append(ans, newInterval)

	// 右边不重叠的，直接入
	for i < size {
		ans = append(ans, intervals[i])
		i++
	}
	return ans
}
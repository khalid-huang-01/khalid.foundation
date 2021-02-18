package main

import "sort"

// 按end排序，查看后续的start是否在end以内，如果是，表明可以一起收拾
func findMinArrowShots(points [][]int) int {
	if len(points) == 0 {
		return 0
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i][1] < points[j][1]
	})
	curX := points[0][1]
	rsl := 1
	for i := 1; i < len(points); i++ {
		if points[i][0] > curX {
			rsl += 1
			curX = points[i][1]
		}
	}
	return rsl
}


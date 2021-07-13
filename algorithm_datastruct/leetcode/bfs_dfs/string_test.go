package bfs_dfs

import "fmt"

//leetcode 93

var (
	segments [4]int
)

func restoreIpAddresses(s string) []string {
	rsl := make([]string, 0)
	backtracking2(s,0, 0, &rsl)
	return rsl
}

func backtracking2(s string, segID int, start int, rsl *[]string) {
	if start == len(s) && segID == 4 {
		tmp := fmt.Sprintf("%d.%d.%d.%d", segments[0], segments[1],segments[2], segments[3])
		*rsl = append(*rsl, tmp)
		return
	}
	if segID >= 4 || start >= len(s) {
		return
	}

	// 遍历当前所有可能情况
	// 前导零只能是一个数字，比如可以是0.0，但不能是00
	// 提前处理前导0，自成一体
	if s[start] == '0' {
		segments[segID] = 0
		backtracking2(s, segID+1, start+1, rsl)
		return
	}
	// 处理其他可以合成一块的情况
	// 回溯不用最后恢复，因为都是副本没影响
	var ipSeg int
	for i := start; i < len(s); i++ {
		ipSeg = 10 * ipSeg + int(s[i]-'0')
		if ipSeg >= 0 && ipSeg <= 255 {
			segments[segID] = ipSeg
			backtracking2(s, segID+1, i+1, rsl)
		} else {
			break
		}
	}
}
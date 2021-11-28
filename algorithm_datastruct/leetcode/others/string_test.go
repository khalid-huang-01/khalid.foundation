package others

// leetcode 777
// 规律题：我们可以看出其实XL->LX是把L左移，RX->XR是把R右移；而且LR之间无法跨越
// 我们可以认为start和end的XR都抽取出来，理应满足：俩个的LR数目一样而且顺序一样
// 但是在同一位置的LR，原start和end里面的下标中，start的R不大于end的R，L不小于end的L下标
// https://leetcode.com/problems/swap-adjacent-in-lr-string/discuss/217070/Python-using-corresponding-position-
type Tuple struct {
	char  rune
	index int
}

func canTransform(start string, end string) bool {
	startLR := make([]Tuple, 0)
	endLR := make([]Tuple, 0)
	for i := 0; i < len(start); i++ {
		if start[i] != 'X' {
			startLR = append(startLR, Tuple{
				char:  rune(start[i]),
				index: i,
			})
		}
		if end[i] != 'X' {
			endLR = append(endLR, Tuple{
				char:  rune(end[i]),
				index: i,
			})
		}
	}
	if len(startLR) != len(endLR) {
		return false
	}
	for i := 0; i < len(startLR); i++ {
		if startLR[i].char != endLR[i].char {
			return false
		}
		if startLR[i].char == 'R' && startLR[i].index > endLR[i].index {
			return false
		}
		if startLR[i].char == 'L' && startLR[i].index < endLR[i].index {
			return false
		}
	}
	return true
}

// leetcode 1930
// 这个题目比较直观，用双指针就可以了
// 固定一个从头开始，然后第二个指针从后面找跟前面等值的，找到了，就看看这俩个指针之间有多少个不同的字符就可以了
func countPalindromicSubsequence(s string) int {
	size := len(s)
	table1 := [26]bool{}
	count := 0
	for i := 0; i < size; i++ {
		if table1[s[i]-'a'] == true {
			continue
		}
		j := size - 1
		table1[s[i]-'a'] = true
		for j > i+1 {
			if s[i] == s[j] {
				break
			}
			j--
		}
		tables2 := [26]bool{}
		for k := i+1; k < j; k++ {
			if tables2[s[k]-'a'] == true {
				continue
			}
			tables2[s[k]-'a'] = true
			count += 1
		}
	}
	return count
}

// leetcode 1529
// 从左往右遍历，如果不匹配就切换一次
func minFlips(target string) int {
	cur := '0'
	nums := 0
	for i := 0; i < len(target); i++ {
		if rune(target[i]) == cur {
			continue
		}
		if cur == '0' {
			cur = '1'
		} else {
			cur = '0'
		}
		nums++

	}
	return nums
}

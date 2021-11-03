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

import "strconv"

// leetcode 228
func summaryRanges(nums []int) []string {
	result := make([]string, 0)
	if len(nums) == 0 {
		return result
	}

	last := nums[0]
	count := 0
	for _, val := range nums[1:] {
		if val == last+count+1 {
			count++
			continue
		}
		if count == 0 {
			result = append(result, strconv.Itoa(last))
		} else {
			result = append(result, strconv.Itoa(last)+"->"+strconv.Itoa(last+count))
		}
		last = val
		count = 0
	}
	return result
}
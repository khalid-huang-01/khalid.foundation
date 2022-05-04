package pointers

import (
	"sort"
	"testing"
)

// leetcode 15
// 这个题目是求三个数，所以我们可以定一个数，然后使用异向双指针找另外两个数字；
// 这个遍历第一个数字的所有可能情况就可以了
func threeSum(nums []int) [][]int {
	ans := make([][]int, 0)
	size,targetSum := len(nums), 0
	left, right := 0, 0
	sort.Ints(nums)
	for i := 0; i < size - 2; i++ {
		targetSum = -nums[i]
		left = i + 1
		right = size - 1
		for left < right {
			if nums[left] + nums[right] == targetSum {
				ans = append(ans, []int{nums[i], nums[left], nums[right]})
				// fmt.Println(nums[i], " ", nums[left], " ", nums[right])
				left += 1
				right -= 1
				// 不能包含相同的组合，把left与right前进到与已经计算过的不一样的数字
				for left < right && nums[left] == nums[left-1] {
					left += 1
				}
				for left < right && nums[right] == nums[right+1] {
					right -= 1
				}
			} else if nums[left] + nums[right] < targetSum {
				left += 1
				for left < right && nums[left] == nums[left-1] {
					left += 1
				}
			} else {
				right -= 1
				for left < right && nums[right] == nums[right+1] {
					right -= 1
				}
			}
		}
		// 相同数字过滤掉，只算一次
		for i < size - 2 && nums[i+1] == nums[i] {
			i+=1
		}
	}
	return ans
}



// 使用双指针来做
// leetcode 88
// 从后面往前面排，merge算法
func merge(nums1 []int, m int, nums2 []int, n int)  {
	for n > 0 {
		if m == 0 || nums2[n-1] >= nums1[m-1] {
			nums1[m+n-1] = nums2[n-1]
			n--
		} else {
			nums1[m+n-1] = nums1[m-1]
			m--
		}
	}
}

// 双指针前进做对比
// 分三种情况：1. 两个数组都还有值，要比较 2. 只有nums1有值了，直接赋值 3. 只有nums2有值了，直接赋值
func merge1(nums1 []int, nums2 []int) []int {
	m, n, cur := 0, 0, 0
	rsl := make([]int, len(nums1)+len(nums2))
	for m < len(nums1) && n < len(nums2) {
		if nums1[m] < nums2[n] {
			rsl[cur] = nums1[m]
			m++
		} else {
			rsl[cur] = nums2[n]
			n++
		}
		cur++
	}
	for m < len(nums1) {
		rsl[cur] = nums1[m]
		m++
		cur++
	}
	for n < len(nums2) {
		rsl[cur] = nums2[n]
		n++
		cur++
	}
	return rsl
}

func merge2(nums1 []int, nums2 []int) []int {
	rsl := make([]int, 0)
	m, n := 0, 0
	for m < len(nums1) && n < len(nums2) {
		if nums1[m] == nums2[n] {
			rsl = append(rsl ,nums1[m])
			m++
			n++
			continue
		}
		if nums1[m] > nums2[n] {
			n++
		} else {
			m++
		}
	}
	return rsl
}

// leetcode 25
type ListNode struct {
    Val int
    Next *ListNode
}

func reverseKGroup(head *ListNode, k int) *ListNode {
	if head == nil {
		return nil
	}

	var pre, cur, next, tail, check *ListNode
	check, tail, cur = head, head, head

	// 判断是否有足够的数量
	var i int
	for i = 0; i < k && check != nil; i++ {
		check = check.Next
	}

	if i != k {
		return head
	}

	for ;i >= 1; i-- {
		next = cur.Next	// 保留后面的元素
		cur.Next = pre // 修改新指向
		pre = cur // 指标推进
		cur = next // 指标推进
	}
	// 每次都把本段的尾巴指向下一段的开头
	tail.Next = reverseKGroup(next, k)
	head = pre // 返回本段的头
	return head
}

func pairWithTargetSum(arr []int, target int) (int, int) {
	left, right := 0, len(arr)-1
	for left < right {
		if arr[left] + arr[right] > target {
			right -= 1
		} else if arr[left] + arr[right] < target {
			left += 1
		} else {
			break
		}
	}
	if arr[left] + arr[right] == target {
		return left, right
	}
	return -1, -1
}

// leetcode 26
// 用同向的快慢指针，left保持左边都是不重复的，如果发现nums[left] != nums[right]，就把nums[right]移动到left右边一个
func removeDuplicates(nums []int) int {
	left, right := 0, 0
	size := len(nums)
	for right < size {
		if nums[left] == nums[right] {
			right += 1
			continue
		}
		nums[left+1] = nums[right]
		left += 1
	}
	return left + 1
}

// leetcode 977
func sortedSquares(nums []int) []int {
	size := len(nums)
	ans := make([]int, size)
	left, right, ansIndex := 0, size-1, size - 1
	leftSquare, rightSquare := 0, 0
	for left < right {
		leftSquare = nums[left] * nums[left]
		rightSquare = nums[right] * nums[right]
		if leftSquare > rightSquare {
			ans[ansIndex] = leftSquare
			left += 1
		} else {
			ans[ansIndex] = rightSquare
			right -= 1
		}
		ansIndex -= 1
	}
	return ans
}

func TestMerge1(t *testing.T) {
	nums1 := []int{1,2,3}
	nums2 := []int{2,5,6}
	t.Log("rsl: ", merge2(nums1, nums2))
}


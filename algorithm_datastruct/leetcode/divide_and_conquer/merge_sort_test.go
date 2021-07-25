package divide_and_conquer

import "testing"


func mergeSort(data []int, left int, right int) {
	// 终止条件
	if left >= right {
		return
	}

	// 分
	middle := left + (right-left)/2
	mergeSort(data, left, middle)
	mergeSort(data, middle+1, right)

	// 治
	merge(data, left, middle, right)
}

func merge(data []int, left, middle, right int) {
	tmp := make([]int, len(data))
	copy(tmp, data)
	//  把data从left到right重新赋值
	k, i, j := left, left, middle+1

	// k为当前data需要赋值的位置
	for k <= right {
		if i > middle {
			// 只有右边可以赋值，左边都已经处理完了
			data[k] = tmp[j]
			j++
		} else if j > right {
			// 只有左边
			data[k] = tmp[i]
			i++
		} else if tmp[i] < tmp[j] {
			// 两边都有，左边小
			data[k] = tmp[i]
			i++
		} else {
			// 两边都有，右边小
			data[k] = tmp[j]
			j++
		}
		k++
	}
}

func TestMergeSort(t *testing.T) {
	data := []int{7, 3, 2, 6, 0, 1, 5, 4}
	mergeSort(data, 0, len(data)-1)
	t.Log(data)
}

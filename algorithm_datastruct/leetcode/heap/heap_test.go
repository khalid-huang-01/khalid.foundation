package heap

import (
	"fmt"
	"testing"
)

// 堆排
func heapSort(arr []int) {
	size := len(arr)
	buildMaxHeap(arr, size)
	for i := size - 1; i > 0; i-- {
		swap(arr, 0, i)
		adjustMaxHeap(arr, 0, i) // 每次都把最后一个给去掉了，也就是把长度减一
	}
}

func buildMaxHeap(arr []int, size int) {
	for i := size / 2 - 1; i >= 0; i-- {
		adjustMaxHeap(arr, i, size)
	}
}

//
func adjustMaxHeap(arr []int, index int, size int) {
	// 找三个里面的最大值
	maxIndex, maxValue := index, arr[index]
	leftIndex := index * 2 + 1
	rightIndex := index * 2 + 2
	if leftIndex < size && maxValue < arr[leftIndex] {
		maxIndex = leftIndex
		maxValue = arr[leftIndex]
	}
	if rightIndex < size && maxValue < arr[rightIndex] {
		maxIndex = rightIndex
		maxValue = arr[rightIndex]
	}
	if maxIndex != index {
		swap(arr, index, maxIndex)
		adjustMaxHeap(arr, maxIndex, size)
	}
}

func swap(arr []int, i int, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func TestHeapSort(t *testing.T) {
	arr := []int{1,4,2,5,3,2}
	heapSort(arr)
	fmt.Println(arr)
}



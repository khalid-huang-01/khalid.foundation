package heap

import (
	"container/heap"
	"fmt"
	"sort"
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
	for i := size/2 - 1; i >= 0; i-- {
		adjustMaxHeap(arr, i, size)
	}
}

//
func adjustMaxHeap(arr []int, index int, size int) {
	// 找三个里面的最大值
	maxIndex, maxValue := index, arr[index]
	leftIndex := index*2 + 1
	rightIndex := index*2 + 2
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
	arr := []int{1, 4, 2, 5, 3, 2}
	heapSort(arr)
	fmt.Println(arr)
}

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func heapSort1() {
	h := &IntHeap{2, 1, 5, 6, 4, 3, 7, 9, 8, 0}
	heap.Init(h)
	heap.Push(h, 6)
	fmt.Println(h)
	a := heap.Pop(h)
	fmt.Println(a)
	fmt.Println(h)
}

func TestHeapSort1(t *testing.T) {
	heapSort1()
}

// leetcode 295
type HeapArr struct {
	sort.IntSlice
}

func (h *HeapArr) Push(v interface{}) {
	h.IntSlice = append(h.IntSlice, v.(int))
}

func (h *HeapArr) Pop() interface{} {
	v := h.IntSlice[h.Len()-1]
	h.IntSlice = h.IntSlice[:h.Len()-1]
	return v
}

type MedianFinder struct {
	minHeap *HeapArr
	maxHeap *HeapArr
}

/** initialize your data structure here. */
func Constructor() MedianFinder {
	return MedianFinder{
		minHeap: new(HeapArr),
		maxHeap: new(HeapArr),
	}
}

func (this *MedianFinder) AddNum(num int)  {
	if this.minHeap.Len() > 0 && num > (*this.minHeap).IntSlice[0] {
		heap.Push(this.minHeap,num)
	} else {
		heap.Push(this.maxHeap,-num)
	}

	if this.minHeap.Len() - this.maxHeap.Len() == 2 {
		heap.Push(this.maxHeap,-(heap.Pop(this.minHeap)).(int))
	} else if this.maxHeap.Len() - this.minHeap.Len() == 2 {
		heap.Push(this.minHeap,-(heap.Pop(this.maxHeap)).(int))
	}
}

func (this *MedianFinder) FindMedian() float64 {
	if this.minHeap.Len() > this.maxHeap.Len() {
		return float64((*this.minHeap).IntSlice[0])
	} else if this.minHeap.Len() < this.maxHeap.Len() {
		return -float64((*this.maxHeap).IntSlice[0])
	}
	return float64((*this.minHeap).IntSlice[0]-(*this.maxHeap).IntSlice[0])/float64(2)
}


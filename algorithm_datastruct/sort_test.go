package main

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

// 每次都把最大的数据冒到坐标最大的位置
func bubbleSort(data []int) {
	size := len(data)
	for i := 0; i < size; i++ { //迭代的次数
		for j := 0; j < size-i-1; j++ { //当前迭代需要比较的次数
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
}

func insertSort(data []int) {
	size := len(data)
	for i := 1; i < size; i++ {
		//在排好序的数组中查找需要这个新数字的位置，每次比较，成功就往后移动
		tmp := data[i] //把i位置的data位置空出来，方便后面存放数据
		var j int
		for j = i - 1; j >= 0; j-- {
			if tmp < data[j] {
				data[j+1] = data[j]
			} else {
				break
			}
		}

		data[j+1] = tmp
	}
}

func _mergeSort(data []int, low int, high int) {
	//合法性检查和终止条件
	if low >= high {
		return
	}
	//分解，归并小解答案
	middle := low + (high-low)/2
	_mergeSort(data, low, middle)
	_mergeSort(data, middle+1, high)
	merge(data, low, middle, high)
}

//赋值的手法高明
func merge(data []int, low int, middle int, high int) {
	tmp := make([]int, len(data))
	copy(tmp, data)
	k, i, j := low, low, middle+1
	//k代表需要为当data赋值的下标
	for k <= high {
		if i > middle {
			data[k] = tmp[j]
			j++
		} else if j > high {
			data[k] = tmp[i]
			i++
		} else if tmp[i] > tmp[j] {
			data[k] = tmp[j]
			j++
		} else {
			data[k] = tmp[i]
			i++
		}
		k++
	}
}

func mergeSort(data []int) {
	low := 0
	high := len(data) - 1
	_mergeSort(data, low, high)
}

func _quickSort(data []int, low int, high int) {
	//终止条件
	if low >= high {
		return
	}
	//划分
	index := partition(data, low, high)

	_quickSort(data, low, index-1)
	_quickSort(data, index+1, high)

}

func partition(data []int, low int, high int) int {
	//随机选择一个数做为基准值，并与high位置交互 => high作为一个基准值了
	swap(data, low+rand.Intn(high-low+1), high)

	//换位,双指针
	var i, j int
	// i是第一个比基准值大的数
	// j是用于判断的数
	for i, j = low, low; j < high; j++ {
		if data[j] < data[high] {
			swap(data, i, j)
			i++
		}
	}

	swap(data, i, high)
	return i

}

func swap(data []int, i int, j int) {
	data[i], data[j] = data[j], data[i]
}

func quickSort(data []int) {
	low := 0
	high := len(data) - 1
	_quickSort(data, low, high)
}

// 堆排序：https://www.cnblogs.com/xingyunshizhe/p/11311754.html

// 使用内置的排序算法 https://www.jianshu.com/p/4bf3c94a15a6
type person struct {
	Name string
	Age  int
}

func TestSort(t *testing.T) {
	// data := []int{4, 1, 2, 6, 8, 3}
	// bubbleSort(data)
	// insertSort(data)
	// mergeSort(data)
	// quickSort(data)
	// fmt.Println(data)
	a := []person{
		{
			Name: "AAA",
			Age:  55,
		},
		{
			Name: "BBB",
			Age:  22,
		},
		{
			Name: "CCC",
			Age:  0,
		},
		{
			Name: "DDD",
			Age:  22,
		},
		{
			Name: "EEE",
			Age:  11,
		},
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i].Age < a[j].Age
	})
	fmt.Println(a)

}

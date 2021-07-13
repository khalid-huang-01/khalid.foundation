package main

import (
	"sort"
	"testing"
)

func TestSortSlice(t *testing.T) {
	array := []int{55,2,32,12,7,3}
	sort.Slice(array, func(i, j int) bool {
		return array[i] < array[j]
	})
	t.Log(array)
}

// 自定义对象的列表
type Person struct {
	Name string
	Age int
}

func TestPersonSliceSort(t *testing.T) {
	persons := []Person{
		{
			Name: "123",
			Age:  20,
		},
		{
			Name: "456",
			Age: 1,
		},
		{
			Name: "789",
			Age: 5,
		},
	}
	sort.Slice(persons, func(i, j int) bool {
		return persons[i].Age < persons[j].Age
	})
	t.Log(persons)
}

// 自定义对象排序
type ByAge []Person

func (b ByAge) Len() int {
	return len(b)
}

func (b ByAge) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByAge) Less(i, j int) bool {
	return b[i].Age < b[j].Age
}

func TestByAgeSort(t *testing.T) {
	persons := []Person{
		{
			Name: "123",
			Age:  20,
		},
		{
			Name: "456",
			Age: 1,
		},
		{
			Name: "789",
			Age: 5,
		},
	}
	sort.Sort(ByAge(persons))
	t.Log(persons)
}

func binarySearch1(arr []int, target int) int {
	high := len(arr)-1
	low := 0
	for low <= high {
		mid := low + (high-low)/2
		if arr[mid] == target {
			return mid
		}
		if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1
}

func binarySearch2(arr []int, target int) {

}

func TestBinarySearch(t *testing.T) {
	array := []int{2,5,6,8,9,11,34,67}
	target := 11
	t.Log(binarySearch1(array, target))
}


package main

arr := []int{2,1,5,61,3,7}
// 数组排序


type Student struct {
	age int
}
// 对象数组排序
func studentListSort() {
	studentList := []Student{
		{age: 1},
		{age: 5},
		{age: 2},
		{age: 10},
		{age: 9},
	}
}

type Students []Student
// 对象列表别名排序
func studentsSort() {
	students := Students{
		{age: 1},
		{age: 5},
		{age: 2},
		{age: 10},
		{age: 9},
	}
}
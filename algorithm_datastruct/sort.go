package main

import "fmt"

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

func main() {
	data := []int{4, 1, 2, 6, 8, 3}
	// bubbleSort(data)
	insertSort(data)
	fmt.Println(data)

}

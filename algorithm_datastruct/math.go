package main

import (
	"fmt"
)

var prime [100]int
var count int

func solution(n int) []int {
	// 迭代，从2到n，对数字进行分解
	var m int
	var result [100]int
	var maxIndex int
	for i := 2; i <= n; i++ {
		m = i
		for j := 0; j < count; j++ {
			for m%prime[j] == 0 {
				m /= prime[j]
				result[j] = result[j] + 1
				if j > maxIndex {
					if j == 3 {
						fmt.Println("result[j]: ", result[j])
					}
					maxIndex = j
				}
			}
		}
	}
	return result[:maxIndex+1]
}

func is_prime(value int) bool {
	for i := 2; i*i <= value; i++ {
		if value%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	// 求出100以内的素数
	for i := 2; i < 100; i++ {
		if is_prime(i) {
			prime[count] = i
			count = count + 1
		}
	}
	// 求出结果
	fmt.Println("prime: ", prime)
	n := 5
	fmt.Println(solution(n))
}

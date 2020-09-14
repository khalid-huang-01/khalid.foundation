package main

import (
	"sort"
	"strings"
)

// 415. Add Strings
// type byte = uint8
func addStrings(num1 string, num2 string) string {
	var result string
	cur1 := len(num1) - 1
	cur2 := len(num2) - 1
	var remain, carry byte

	for cur1 >= 0 || cur2 >= 0 {
		tmp := carry
		carry = 0
		if cur1 >= 0 {
			tmp += num1[cur1] - '0'
			cur1--
		}
		if cur2 >= 0 {
			tmp += num2[cur2] - '0'
			cur2--
		}
		remain = tmp % 10
		carry = tmp / 10
		result = string(remain+'0') + result
	}
	if carry > 0 {
		result = string(carry+'0') + result
	}
	return result
}

// leetcode 820
// 这个题目可以使用排序也可以使用字典树
// 主要目标是为了查出是否有单词是另一个单词的后缀（前缀反过来就可以了），然后排除掉是的，对剩下的进行计算
func minimumLengthEncoding(words []string) int {
//	进行反转
	for i := 0; i < len(words); i++ {
		words[i] = reverseStrings(words[i])
	}
//  进行升序排序
	sort.Strings(words)
	words = append(words, "*") //保证全部遍历
	res := 0
//  进行判断是否是前缀
	for i := 0; i < len(words) - 1; i++ {
		if strings.HasPrefix(words[i+1], words[i]) {
			continue
		}
		res = res + len(words[i]) + 1
	}
	return res
}

// 需要先转成rune，底层是转成了byte[]
func reverseStrings(s string) string {
	//使用runes
	//runes := []rune(s)
	//for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
	//	runes[from], runes[to] = runes[to], runes[from]
	//}
	//return string(runes)
	//使用bytes，一般来说string的char操作，转成byte好操作
	bytes := []byte(s)
	for from, to := 0, len(bytes) - 1; from < to; from, to = from + 1, to -1 {
		bytes[from], bytes[to] = bytes[to], bytes[from]
	}
	return string(bytes)
}


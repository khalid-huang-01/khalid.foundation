func minimumLengthEncoding(words []string) int {
	//反转
	for i, word := range(words) {
		words[i] = reverseString(word)
	}
	//安排字典序排序
	sort.Strings(word)
	//判断是否是前缀，进行处理
	var result int
	for i := 0; i < len(word) - 1; i++ {
		if strings.HasPrefix(word[i+1], word[i]) {
			continue
		}
		result = result + len(word[i]) + 1
	}
	return res
}

func reverseString(word string) {
	//转成rune
	word := rune(word)
	for from, to := 0, len(word) - 1; from < to; from, to = from + 1, to -1- {
		word[from], word[to] = word[to], word[from]
	}
	return string(word)
}
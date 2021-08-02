package leetcode

// 前缀树思想

type TrieNode struct {
	ChildrenCount int
	children [26]*TrieNode
	Used bool // 用于统计的时候表示是否被用过，用于避免由多个一一样的单词
}
func minimumLengthEncoding(words []string) int {
	root := &TrieNode{}
	wordFirstLetterNode := make([]*TrieNode, len(words))
	for index, word := range words {
		cur := root
		for i := len(word) - 1; i >= 0; i-- {
			cur = getOrInsert(cur, word[i])
		}
		wordFirstLetterNode[index] = cur
	}
	// 通过统计每个单词的第一个字母所对应的后续字母是不是没有为0（要倒过来的），就知道这个是不是没法被其他单词包含的
	rsl := 0
	for index, node := range wordFirstLetterNode {
		if node.ChildrenCount == 0 && !node.Used {
			rsl += len(words[index]) + 1
			node.Used = true
		}
	}
	return rsl
}

func getOrInsert(parent *TrieNode, letter byte) *TrieNode {
	i := letter - 'a'
	if parent.children[i] == nil {
		parent.children[i] = &TrieNode{
			ChildrenCount: 0,
			Used: false,
		}
		parent.ChildrenCount += 1
	}
	return parent.children[i]
}
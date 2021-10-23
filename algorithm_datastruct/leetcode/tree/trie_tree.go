package tree

// 前缀树思想

// leetcode 208
type Trie struct {
	// 不需要有value，因为最顶部就是一个空的
	children [26]*Trie
	isEnd bool // 表示是不是一个单词的结尾
}


func Constructor() Trie {
	return Trie{}
}

// this就是根节点
func (this *Trie) Insert(word string)  {
	cur := this
	for _, letter := range word {
		if cur.children[letter - 'a'] == nil {
			cur.children[letter - 'a'] = &Trie{}
		}
		cur = cur.children[letter - 'a']
	}
	cur.isEnd = true
}


func (this *Trie) Search(word string) bool {
	cur := this
	for _, letter := range word {
		if cur.children[letter - 'a'] == nil {
			return false
		}
		cur = cur.children[letter - 'a']
	}
	return cur.isEnd
}


func (this *Trie) StartsWith(prefix string) bool {
	cur := this
	for _, letter := range prefix {
		if cur.children[letter - 'a'] == nil {
			return false
		}
		cur = cur.children[letter - 'a']
	}
	return true
}

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

// leetcode 211
type WordDictionary struct {
	children [26]*WordDictionary
	isEnd bool
}


func Constructor1() WordDictionary {
	return WordDictionary{}
}


func (this *WordDictionary) AddWord(word string)  {
	cur := this
	for _, letter := range word {
		if cur.children[letter - 'a'] == nil {
			cur.children[letter - 'a'] = &WordDictionary{}
		}
		cur = cur.children[letter - 'a']
	}
	cur.isEnd = true
}


func (this *WordDictionary) Search(word string) bool {
	cur := this
	for i, letter := range word {
		if letter == '.' {
			// 这里要用回溯法，不断去试每个child匹配.后，后续是否符合情况； 这里的回溯没有状态需要恢复
			for _, child := range cur.children {
				if child == nil {
					continue
				}
				// 使用当前child来match.，继续往下走
				if child.Search(word[i+1:]) {
					return true
				}
				// 继续，使用后面的child来match
			}
			// 如果.没法匹配上，就返回false
			return false
		} else {
			if cur.children[letter - 'a'] == nil {
				return false
			}
			cur = cur.children[letter - 'a']
		}
	}
	return cur.isEnd
}


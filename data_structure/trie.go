package data_structure

type Trie struct {
	count int64
	root  *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{},
	}
}

type TrieNode struct {
	count   int64
	prefix  int64
	indexes []int64
	next    [26]*TrieNode
}

func (trie *Trie) Insert(words []string) {
	for _, word := range words {
		root := trie.root
		for i := range word {
			if root.next[word[i]-'a'] == nil {
				root.next[word[i]-'a'] = &TrieNode{}
			}
			root = root.next[word[i]-'a']
			root.prefix++
		}
		root.count++
		root.indexes = append(root.indexes, trie.count)
		trie.count++
	}
}

func (trie *Trie) Search(word string) (int64, []int64) {
	root := trie.root
	for i := range word {
		root = root.next[word[i]-'a']
		if root == nil {
			return 0, nil
		}
	}
	return root.count, root.indexes
}

func (trie *Trie) SearchPrefix(prefix string) int64 {
	root := trie.root
	for i := range prefix {
		root = root.next[prefix[i]-'a']
		if root == nil {
			return 0
		}
	}
	return root.prefix
}

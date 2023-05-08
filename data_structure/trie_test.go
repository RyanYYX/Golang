package data_structure

import (
	"fmt"
	"strings"
	"testing"
)

var trie *Trie

func init() {
	trie = NewTrie()
	trie.Insert(strings.Split("i love you this world ice lo yy i you i you", " "))
}

func TestTrie_Search(t *testing.T) {
	count, indexes := trie.Search("you")
	fmt.Printf("%d, %v", count, indexes)
}

func TestTrie_SearchPrefix(t *testing.T) {
	count := trie.SearchPrefix("y")
	fmt.Printf("%d", count)
}

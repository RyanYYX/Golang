package data_structure

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	_maxLevel uint8 = 32
)

type SkipList struct {
	level int
	head  *SkipListNode
}

func NewSkipList() *SkipList {
	return new(SkipList)
}

type SkipListNode struct {
	score int64
	member string
	next  []*SkipListNode
}

func (node *SkipListNode) String() string {
	return fmt.Sprintf("score: %d, member: %s\n", node.score, node.member)
}

func (list *SkipList) getRandomLevel() uint8 {
	rand.Seed(time.Now().UnixNano())
	random := rand.Uint32()
	for i := uint8(0); i < _maxLevel; i++ {
		if random&(1<<(_maxLevel-1-i)) == 1 {
			return i
		}
	}
	return 0
}

func (list *SkipList) Insert(score int64, member string) bool {
	level := list.getRandomLevel()
	if list.head == nil {
		list.head = &SkipListNode{}
		list.level = -1
	}

	node := &SkipListNode{score: score, member: member}
	node.next = make([]*SkipListNode, level+1)
	current := list.head
	for i := list.level; i >= 0; i-- {
		next := current.next[i]
		for next != nil && next.score < score {
			current, next = next, next.next[i]
		}

		if i <= int(level) {
			node.next[i] = current.next[i]
			current.next[i] = node
		}
	}

	if int(level) > list.level {
		for i := list.level; i < int(level); i++ {
			list.head.next = append(list.head.next, node)
		}
		list.level = int(level)
	}
	return true
}

func (list *SkipList) Delete(score int64) bool {
	current := list.head
	for i := list.level; i >= 0; i-- {
		next := current.next[i]
		for next != nil && next.score < score {
			current = next
			next = current.next[i]
		}

		if next != nil && next.score == score {
			current.next[i] = next.next[i]
			next.next[i] = nil
		}
	}
	return true
}

func (list *SkipList) Query(score int64) *SkipListNode {
	current := list.head
	for i := list.level; i >= 0; i-- {
		next := current.next[i]
		for next != nil && next.score < score {
			current = next
			next = current.next[i]
		}

		if next != nil && next.score == score {
			return next
		}
	}
	return nil
}

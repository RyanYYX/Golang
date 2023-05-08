package data_structure

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSkipList(t *testing.T) {
	list := NewSkipList()
	for i := 0; i < 50000; i++ {
		random := rand.Int63n(50000)
		a := time.Now().UnixNano()
		list.Insert(random, fmt.Sprint(random))
		b := time.Now().UnixNano()
		print((b-a)/1000, " us", "\n")
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		node := list.Query(rand.Int63n(50000))
		if node != nil {
			print(node.String())
		}
	}
}

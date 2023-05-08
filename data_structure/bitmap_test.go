package data_structure

import (
	"fmt"
	"testing"
)

func TestBitMap(t *testing.T) {
	bitmap := NewBitMap()
	bitmap.Add(0)
	bitmap.Add(1)
	bitmap.Add(2)
	bitmap.Add(3)
	fmt.Printf("%+v, %d, %v, %v, %v, %v\n", bitmap.ToArray(), bitmap.Count(),
		bitmap.Has(0), bitmap.Has(1), bitmap.Has(2), bitmap.Has(3))
	bitmap.Del(3)
	bitmap.Del(2)
	bitmap.Del(1)
	fmt.Printf("%+v, %d, %v, %v, %v, %v\n", bitmap.ToArray(), bitmap.Count(),
		bitmap.Has(0), bitmap.Has(1), bitmap.Has(2), bitmap.Has(3))
}

func TestRoaringBitMap(t *testing.T) {
	bitmap := NewRoaringBitMap()
	bitmap.Add(0)
	bitmap.Add(1)
	bitmap.Add(2)
	bitmap.Add(3)
	fmt.Printf("%v, %v, %v, %v\n",
		bitmap.Has(0), bitmap.Has(1), bitmap.Has(2), bitmap.Has(3))
	bitmap.Del(3)
	bitmap.Del(2)
	bitmap.Del(1)
	fmt.Printf("%v, %v, %v, %v\n",
		bitmap.Has(0), bitmap.Has(1), bitmap.Has(2), bitmap.Has(3))
}

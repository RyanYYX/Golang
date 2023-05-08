package data_structure

type BitMap struct {
	count int
	words []uint64
}

func NewBitMap() *BitMap {
	return new(BitMap)
}

func (bm *BitMap) Has(x uint64) bool {
	word, bit := x/64, x%64
	if word >= uint64(len(bm.words)) {
		bm.words = append(bm.words, 0)
	}
	return (bm.words[word] >> bit) & 1 == 1
}

func (bm *BitMap) Add(x uint64) bool {
	word, bit := x/64, x%64
	if word >= uint64(len(bm.words)) {
		bm.words = append(bm.words, 0)
	}

	if bm.words[word] >> bit & 1 == 1 {
		return false
	}

	bm.words[word] |= 1 << bit
	bm.count++
	return true
}

func (bm *BitMap) Del(x uint64) bool {
	word, bit := x/64, x%64
	if word >= uint64(len(bm.words)) {
		bm.words = append(bm.words, 0)
	}

	if bm.words[word] >> bit & 1 != 1 {
		return false
	}
	bm.words[word] &= ^(1 << bit)
	bm.count--
	return true
}

func (bm *BitMap) ToArray() []uint64 {
	array, i := make([]uint64, bm.count), 0
	for word, offset := range bm.words {
		for bit := 0; offset > 0; offset, bit = offset >> 1, bit + 1 {
			if offset & 1 == 1 {
				array[i] = uint64(bit + word*64)
				i = i + 1
			}
		}
	}
	return array
}

func (bm *BitMap) Count() int {
	return bm.count
}
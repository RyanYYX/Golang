package data_structure

type status uint8

const (
	_statusArray status = iota
	_statusBitmap

	_maxArrayLength   = 4096
	_toArrayThreshold = _maxArrayLength * 0.75
)

type RoaringBitMap struct {
	status  map[uint16]status
	buckets map[uint16]interface{}
}

func NewRoaringBitMap() *RoaringBitMap {
	return &RoaringBitMap{
		status:  make(map[uint16]status),
		buckets: make(map[uint16]interface{}),
	}
}

func (rbm *RoaringBitMap) Has(x uint32) bool {
	high, low := uint16(x>>16), uint16((x<<16)>>16)
	switch rbm.status[high] {
	case _statusArray:
		if array, ok := rbm.buckets[high].([]uint16); ok {
			_, hit := rbm.binarySearch(array, low)
			return hit
		}
	case _statusBitmap:
		if bitmap, ok := rbm.buckets[high].(*BitMap); ok {
			return bitmap.Has(uint64(low))
		}
	}
	return false
}

func (rbm *RoaringBitMap) Add(x uint32) bool {
	high, low := uint16(x>>16), uint16((x<<16)>>16)
	if _, ok := rbm.buckets[high]; !ok {
		array := make([]uint16, 1)
		array[0] = low
		rbm.buckets[high] = array
		rbm.status[high] = _statusArray
		return true
	}

	switch rbm.status[high] {
	case _statusArray:
		return rbm.arrayAdd(high, low)
	case _statusBitmap:
		return rbm.bitmapAdd(high, low)
	}
	return false
}

func (rbm *RoaringBitMap) arrayAdd(high, low uint16) bool {
	array := rbm.buckets[high].([]uint16)
	index, hit := rbm.binarySearch(array, low)
	if hit {
		return false
	}
	array = append(array, 0)
	copy(array[index+1:], array[index:])
	array[index] = low
	rbm.buckets[high] = array
	if len(array) >= _maxArrayLength {
		bitmap := rbm.arrayToBitmap(array)
		rbm.buckets[high] = bitmap
		rbm.status[high] = _statusBitmap
	}
	return true
}

func (rbm *RoaringBitMap) bitmapAdd(high, low uint16) bool {
	bitmap := rbm.buckets[high].(*BitMap)
	return bitmap.Add(uint64(low))
}

func (rbm *RoaringBitMap) Del(x uint32) bool {
	high, low := uint16(x>>16), uint16((x<<16)>>16)
	switch rbm.status[high] {
	case _statusArray:
		return rbm.arrayDel(high, low)
	case _statusBitmap:
		return rbm.bitmapDel(high, low)
	}
	return false
}

func (rbm *RoaringBitMap) arrayDel(high, low uint16) bool {
	array := rbm.buckets[high].([]uint16)
	index, hit := rbm.binarySearch(array, low)
	if !hit {
		return false
	}
	copy(array[index:], array[index+1:])
	rbm.buckets[high] = array[:len(array)-1]
	return true
}

func (rbm *RoaringBitMap) bitmapDel(high, low uint16) bool {
	bitmap := rbm.buckets[high].(*BitMap)
	success := bitmap.Del(uint64(low))
	if !success {
		return false
	}
	if bitmap.Count() < _toArrayThreshold {
		array := rbm.bitmapToArray(bitmap)
		rbm.buckets[high] = array
		rbm.status[high] = _statusArray
	}
	return true
}

func (rbm *RoaringBitMap) arrayToBitmap(array []uint16) *BitMap {
	bitmap := NewBitMap()
	for _, num := range array {
		bitmap.Add(uint64(num))
	}
	return bitmap
}

func (rbm *RoaringBitMap) bitmapToArray(bitmap *BitMap) []uint16 {
	array := make([]uint16, bitmap.Count())
	for _, num := range bitmap.ToArray() {
		array = append(array, uint16(num))
	}
	return array
}

func (rbm *RoaringBitMap) binarySearch(array []uint16, low uint16) (index int, hit bool) {
	l, r := 0, len(array)
	for l < r {
		mid := (l + r) / 2
		if array[mid] == low {
			return mid, true
		} else if array[mid] < low {
			l = mid + 1
		} else {
			r = mid
		}
	}
	return l, false
}

package algorism

import (
	"fmt"
	"testing"
)

func TestSelection(t *testing.T) {
	nums := []int{30, 20, 40, 10, 50, 1, 5, 9, 41, 67, 31, 77}
	selectionSort(nums)
	fmt.Print(nums)
}

func TestBubble(t *testing.T) {
	nums := []int{30, 20, 40, 10, 50, 1, 5, 9, 41, 67, 31, 77}
	bubbleSort(nums)
	fmt.Print(nums)
}

func TestInsertion(t *testing.T) {
	nums := []int{30, 20, 40, 10, 50, 1, 5, 9, 41, 67, 31, 77}
	insertionSort(nums)
	fmt.Print(nums)
}

func TestShell(t *testing.T) {
	nums := []int{30, 20, 40, 10, 50, 1, 5, 9, 41, 67, 31, 77}
	shellSort(nums)
	fmt.Print(nums)
}

func TestQuick(t *testing.T) {
	nums := []int{30, 20, 40, 10, 50, 1, 5, 9, 41, 67, 31, 77}
	quickSort(nums, 0, len(nums)-1)
	fmt.Print(nums)
}

func TestHeap(t *testing.T) {
	nums := []int{30, 20, 40, 10, 50, 1, 5, 9, 41, 67, 31, 77}
	heapSort(nums)
	fmt.Print(nums)
}

func TestMerge(t *testing.T) {
	nums := []int{30, 20, 40, 10, 50, 1, 5, 9, 41, 67, 31, 77}
	mergeSort(nums, 0, len(nums)-1)
	fmt.Print(nums)
}

func TestBucket(t *testing.T) {
	nums := []int{30, 20, 40, 10, 50, 1, 5, 9, 41, 67, 31, 77}
	bucketSort(nums)
	fmt.Print(nums)
}

func TestRadix(t *testing.T) {
	nums := []int{30, 20, 40, 10, 50, 1, 5, 9, 41, 67, 31, 77}
	radixSort(nums)
	fmt.Print(nums)
}

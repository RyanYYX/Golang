package algorism

// 选择排序
func selectionSort(nums []int) {
	for i:=0; i<len(nums)-1; i++ {
		var min = i
		for j:=i+1; j<len(nums); j++ {
			if nums[min] >= nums[j] {
				min = j
			}
		}
		nums[i], nums[min] = nums[min], nums[i]
	}
}

// 冒泡排序
func bubbleSort(nums []int) {
	for i:=0; i<len(nums)-1; i++ {
		for j:=len(nums)-1; j>i; j-- {
			if nums[j] < nums[j-1] {
				nums[j], nums[j-1] = nums[j-1], nums[j]
			}
		}
	}
}

// 插入排序
func insertionSort(nums []int) {
	for i:=1; i<len(nums); i++ {
		for j:=i-1; j>=0; j-- {
			if nums[j] > nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
			} else {
				break
			}
		}
	}
}

// 希尔排序
func shellSort(nums []int) {
	for i:=len(nums)/2; i>0; i/=2 {
		for j:=i; j<len(nums)-1; j++ {
			for k:=j-i; k>=0; k-=i {
				if nums[k] > nums[k+i] {
					nums[k], nums[k+i] = nums[k+i], nums[k]
				}
			}
		}
	}
}

func quick(nums []int, l, r int) int {
	var temp = nums[l]
	for ; l<r; {
		for ; l<r && nums[r]>temp; r-- {}
		nums[l] = nums[r]
		for ; l<r && nums[l]<temp; l++ {}
		nums[r] = nums[l]
	}
	nums[l] = temp
	return l
}

// 快速排序
func quickSort(nums []int, l, r int) {
	if l < r {
		k := quick(nums, l, r)
		quickSort(nums, l, k-1)
		quickSort(nums, k+1, r)
	}
}

func maximumHeap(nums []int, i, length int) {
	var flag bool
	for j:=2*i+1; !flag && j<length; i, j = j, 2*j+1  {
		if j+1<length && nums[j] < nums[j+1] {
			j++
		}

		if nums[i] < nums[j] {
			nums[i], nums[j] = nums[j], nums[i]
		} else {
			flag = true
		}
	}
}

// 堆排序
func heapSort(nums []int) {
	for i:=(len(nums)+1)/2-1; i >= 0; i-- {
		maximumHeap(nums, i, len(nums))
	}
	for i:=len(nums)-1; i>0; i-- {
		nums[0], nums[i] = nums[i], nums[0]
		maximumHeap(nums, 0, i)
	}
}

// 归并排序
func mergeSort(nums []int, l, r int) {
	if l == r {
		return
	}

	middle := (l + r) / 2
	mergeSort(nums, l, middle)
	mergeSort(nums, middle+1, r)
	var temp []int
	i, j := l, middle+1
	for ; i <= middle && j <= r; {
		if nums[i] < nums[j] {
			temp = append(temp, nums[i])
			i++
		} else {
			temp = append(temp, nums[j])
			j++
		}
	}

	for ; i <= middle; i++ {
		temp = append(temp, nums[i])
	}

	for ; j <= r; j++ {
		temp = append(temp, nums[j])
	}

	for i:=0; l <= r; l, i = l+1, i+1 {
		nums[l] = temp[i]
	}
}

// 桶排序，范围：[0, 100)
func bucketSort(nums []int) {
	var buckets [][]int
	buckets = make([][]int, 10)
	for _, num := range nums {
		idx := num / 10
		buckets[idx] = append(buckets[idx], num)
		bucket := buckets[idx]
		for i:=len(bucket)-2; i>=0; i++ {
			if bucket[i] > bucket[i+1] {
				bucket[i], bucket[i+1] = bucket[i+1], bucket[i]
			} else {
				break
			}
		}
	}

	var temp []int
	for _, bucket := range buckets {
		temp = append(temp, bucket...)
	}

	for i:=0; i<len(nums); i++ {
		nums[i] = temp[i]
	}
}

// 基数排序
func radixSort(nums []int) {
	var (
		hex int = 1
		buckets [][]int
	)

	for {
		buckets = make([][]int, 10)
		for _, num := range nums {
			idx := num / hex % 10
			buckets[idx] = append(buckets[idx], num)
		}

		if len(buckets[0]) == len(nums) {
			break
		}

		var temp []int
		for _, bucket := range buckets {
			temp = append(temp, bucket...)
		}

		for i := range nums {
			nums[i] = temp[i]
		}
		hex *= 10
	}
}
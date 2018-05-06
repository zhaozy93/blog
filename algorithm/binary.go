package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	length := 20
	// array := genRandomArr(length)
	// value := array[rand.Intn(length)]

	array := []int{38, 104, 156, 246, 291, 368, 410, 507, 580, 613, 651, 660, 753, 806, 883, 935, 1030, 1059, 1060, 1088}
	value := 1089

	fmt.Printf("Generate Array success, %v, %d\n", array, value)
	index := binarySearch(array, value, 0, length-1)
	fmt.Println(index)
}

func binarySearch(arr []int, value, start, end int) int {
	if start > end {
		return -1
	}
	index := -1
	half := (end-start)/2 + start
	if arr[half] == value {
		return half
	}
	if arr[half] < value {
		return binarySearch(arr, value, half+1, end)
	}
	if arr[half] > value {
		return binarySearch(arr, value, start, half-1)
	}
	return index
}

func genRandomArr(n int) []int {
	arr := make([]int, n, n)
	rand.Seed(time.Now().Unix())
	base := rand.Intn(100)
	arr[0] = base
	for i := 1; i < n; i++ {
		rand.Seed(time.Now().UnixNano())
		offset := rand.Intn(100)
		arr[i] = arr[i-1] + offset
	}
	return arr
}

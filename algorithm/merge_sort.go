package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	length := 20
	array := genRandomArr(length)
	fmt.Printf("Generate Array success, %v\n", array)
	array = mergeSort(array, 0, length-1)
	fmt.Println(array)
}

func mergeSort(arr []int, start, end int) []int {
	half := (end + start) / 2
	if end-start == 1 {
		if arr[start] < arr[end] {
			return []int{arr[start], arr[end]}
		}
		return []int{arr[end], arr[start]}
	}
	if end-start == 0 {
		return []int{arr[start]}
	}
	return sort(mergeSort(arr, start, half), mergeSort(arr, half+1, end))
}

func sort(arr1 []int, arr2 []int) []int {
	fmt.Println("merge", arr1, arr2)
	n1, n2 := len(arr1), len(arr2)
	sum := n1 + n2
	arr := make([]int, sum, sum)
	for i, i1, i2 := 0, 0, 0; i < sum; i++ {
		if i1 == n1 {
			arr[i] = arr2[i2]
			i2++
			continue
		}
		if i2 == n2 {
			arr[i] = arr1[i1]
			i1++
			continue
		}
		if arr1[i1] < arr2[i2] {
			arr[i] = arr1[i1]
			i1++
			continue
		} else {
			arr[i] = arr2[i2]
			i2++
			continue
		}
	}
	return arr
}

func genRandomArr(n int) []int {
	arr := make([]int, n, n)
	rand.Seed(time.Now().Unix())
	base := rand.Intn(100)
	arr[0] = base
	for i := 1; i < n; i++ {
		rand.Seed(time.Now().UnixNano())
		offset := rand.Intn(100)
		if offset >= 50 {
			arr[i] = arr[i-1] + offset
		} else {
			arr[i] = arr[i-1] - offset
		}
	}
	return arr
}

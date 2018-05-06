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
	quickSort(array, 0, length-1)
	fmt.Println(array)
}

func quickSort(arr []int, start, end int) {
	// fmt.Println(arr, start, end)
	if start >= end {
		return
	}
	q := partition(arr, start, end)
	fmt.Println(arr, q)
	quickSort(arr, start, q-1)
	quickSort(arr, q+1, end)

}

func partition(arr []int, start, end int) int {
	val := arr[end]
	index := start
	for i := start; i < end; i++ {
		if arr[i] < val {
			arr[i], arr[index] = arr[index], arr[i]
			index++
		}
	}
	arr[index], arr[end] = arr[end], arr[index]
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
		if offset >= 50 {
			arr[i] = arr[i-1] + offset
		} else {
			arr[i] = arr[i-1] - offset
		}
	}
	return arr
}

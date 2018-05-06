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
	array = selectionSort(array)
	fmt.Println(array)
}

func selectionSort(arr []int) []int {
	length := len(arr)
	for i := 0; i < length; i++ {
		smallest := arr[i]
		smallestIndex := i
		for j := i; j < length; j++ {
			if arr[j] < smallest {
				smallest = arr[j]
				smallestIndex = j
			}
		}
		if i != smallestIndex {
			arr[smallestIndex] = arr[i]
			arr[i] = smallest
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

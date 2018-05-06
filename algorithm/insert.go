package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	length := 20
	array := genRandomArr(length)
	// array[0] = 200000
	// array[1] = 30
	fmt.Printf("Generate Array success, %v\n", array)
	array = insertSort(array)
	fmt.Println(array)
}

func insertSort1(arr []int) []int {
	for i := 1; i < len(arr); i++ {
		for j := i; j > 0; j-- {
			if arr[j] < arr[j-1] {
				arr[j-1], arr[j] = arr[j], arr[j-1]
			}
		}
	}
	return arr
}

func insertSort(arr []int) []int {
	for i := 1; i < len(arr); i++ {
		value := arr[i]
		j := i - 1
		for j >= 0 && value < arr[j] {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = value
		fmt.Println(arr)

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

package main

import (
	"fmt"
)

var n = 10
var origin = make([]int, n)

func main() {
	revisePorker(n, 1, 0)
	fmt.Println(origin)
}

func revisePorker(n, current, index int) {
	// fmt.Println("currnt is ", current, "index is ", index, "origin is ", origin, "next is ", next)
	origin[index] = current
	if current == n {
		return
	}
	revisePorker(n, current+1, findNext(n, current, index))
}

func findNext(n, current, index int) int {
	count := 0
	for {
		index++
		if index == n {
			index = 0
		}
		if origin[index] == 0 {
			count++
		}
		if count == 2 {
			return index
		}
	}
}

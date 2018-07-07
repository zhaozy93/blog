package main

import (
	"fmt"
	"math"
)

// [1 6 2 10 3 7 4 9 5 8]
var n = 5
var origin = make([]int, n)

func main() {
	revisePorker(1, 1, 1)
	fmt.Println(origin)
}

func revisePorker(current, index, loop int) {
	// fmt.Println("currnt is ", current, "index is ", index, "origin is ", origin, "loop is ", loop)
	origin[index-1] = current
	if current == n {
		return
	}
	next, loop := findNext(current, index, loop)
	revisePorker(current+1, next, loop)
}

func findNext(current, index, loop int) (int, int) {
	step := int(math.Pow(2, float64(loop)))
	if index+step > n {
		right := 0
		for i := index; i <= n; i++ {
			if origin[i-1] == 0 {
				right++
			}

		}
		step = step * 2
		return (index + step - right*step/2 + right) % n, loop + 1
	} else {
		return index + step, loop
	}
}

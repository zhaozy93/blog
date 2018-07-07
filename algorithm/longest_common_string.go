package main

import (
	"fmt"
)

func main() {
	x := "CATCGA"
	y := "GTACCGTCA"
	fmt.Println(x)
	fmt.Println(y)
	table := computeLcsTable(x, y)
	fmt.Println(table)
	lcs := make([]byte, 4)
	assembleLcs(x, y, table, len(x)-1, len(y)-1, lcs)
	fmt.Println(lcs)
	fmt.Println(string(lcs))
}

func computeLcsTable(x, y string) [][]int {
	xs := []byte(x)
	ys := []byte(y)

	res := make([][]int, len(xs))
	// fmt.Println(res)
	for i := 0; i < len(xs); i++ {
		res[i] = make([]int, len(ys))
		for j := 0; j < len(ys); j++ {
			if j == 0 || i == 0 {
				res[i][j] = 0
				continue
			}
			if xs[i] == ys[j] {
				res[i][j] = res[i-1][j-1] + 1
			} else if res[i][j-1] > res[i-1][j] {
				res[i][j] = res[i][j-1]
			} else {
				res[i][j] = res[i-1][j]
			}
		}
	}
	return res
}

func assembleLcs(x, y string, l [][]int, i, j int, lcs []byte) {
	if l[i][j] == 0 {
		return
	}
	fmt.Printf("current i, j is %d, %d, current lcsLen is %d, lcs is %v\n", i, j, l[i][j], lcs)
	if x[i] == y[j] {
		fmt.Printf("x %d ,y %d is equal\n", i, j)
		lcs[l[i][j]-1] = x[i]
		assembleLcs(x, y, l, i-1, j-1, lcs)
	} else if l[i-1][j] == l[i][j-1] {
		fmt.Printf("x %d ,y %d both pre is equal\n", i, j)

		assembleLcs(x, y, l, i-1, j, lcs)
	} else if l[i-1][j] > l[i][j-1] {
		fmt.Printf("x %d  pre is longer than y %d\n", i, j)

		assembleLcs(x, y, l, i-1, j, lcs)
	} else {
		fmt.Printf("x %d pre is shorter than y %d\n", i, j)

		assembleLcs(x, y, l, i, j-1, lcs)
	}
}

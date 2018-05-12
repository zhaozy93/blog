/*
	有向图 求最短路径  Floyd 算法
*/
package main

import (
	// "errors"
	"fmt"
	"math"
)

// 所有路径都是双向的
var route = map[string]int{
	"aa": 0,
	"bb": 0,
	"cc": 0,
	"dd": 0,
	"ee": 0,
	"ab": 1,
	"ac": 4,
	"be": 2,
	"bc": 3,
	"bd": 2,
	"dc": 5,
	"db": 1,
	"ed": 3,
}

var dist = map[string]int{
	"a": 0,
	"b": math.MaxInt32,
	"c": math.MaxInt32,
	"d": math.MaxInt32,
	"e": math.MaxInt32,
}
var res = map[string]map[string]int{
	"a": map[string]int{
		"a": 0,
		"b": 2,
		"c": 4,
		"d": math.MaxInt32,
		"e": math.MaxInt32,
	},
	"b": map[string]int{
		"a": math.MaxInt32,
		"b": 0,
		"c": 3,
		"d": 2,
		"e": 2,
	},
	"c": map[string]int{
		"a": math.MaxInt32,
		"b": math.MaxInt32,
		"c": 0,
		"d": math.MaxInt32,
		"e": math.MaxInt32,
	},
	"d": map[string]int{
		"a": math.MaxInt32,
		"b": 1,
		"c": 5,
		"d": 0,
		"e": math.MaxInt32,
	},
	"e": map[string]int{
		"a": math.MaxInt32,
		"b": math.MaxInt32,
		"c": math.MaxInt32,
		"d": 3,
		"e": 0,
	},
}

func main() {
	Floyd()
	fmt.Println(res)
}

func Floyd() {
	for port, _ := range dist {
		for start, _ := range dist {
			fmt.Printf("current end is: %s\n", end)
			if start == port {
				continue
			}
			for end, _ := range dist {
				if port == end || start == end {
					continue
				}
				fmt.Printf("current port is: %s\n", port)
				fmt.Printf("%d, %d\n", findShorestDist(start, port), findShorestDist(port, end))
				if findShorestDist(start, port)+findShorestDist(port, end) < res[start][end] {
					fmt.Printf("shortest is: %d\n", findShorestDist(start, port)+findShorestDist(port, end))
					res[start][end] = findShorestDist(start, port) + findShorestDist(port, end)
				}
			}

		}

		if start == "a" {
			return
		}
	}
}

func findShorestDist(start, end string) int {
	val := math.MaxInt32
	if value, ok := route[start+end]; ok {
		val = value
	}
	if val < res[start][end] {
		return val
	}
	return res[start][end]
}

// func relax(uv string) {
// 	start := string(uv[0])
// 	end := string(uv[1])
// 	weight, ok := route[uv]
// 	if !ok {
// 		return
// 	}
// 	if dist[start]+weight < dist[end] {
// 		dist[end] = dist[start] + weight
// 		pred[end] = start
// 	}
// }

func findConnPoint(v string) []string {
	next := make([]string, 0, 0)
	for key, _ := range route {
		if key[0] == ([]byte(v))[0] {
			next = append(next, string(key[1]))
		}
		if key[1] == ([]byte(v))[0] {
			next = append(next, string(key[0]))
		}
	}
	return next
}

func findRouteWeight(u, v string) int {
	if weight, ok := route[u+v]; ok {
		return weight
	}
	return math.MaxInt32
}

func isInLine(v string, line []string) bool {
	for _, val := range line {
		if val == v {
			return true
		}
	}
	return false
}

// func findNextPoint() (string, int, error) {
// 	if len(rest) == 0 {
// 		return "", -1, errors.New("no point in rest arr")
// 	}
// 	next := rest[0]
// 	index := 0
// 	for i, key := range rest {
// 		if weight[key] < weight[next] {
// 			next = key
// 			index = i
// 		}
// 	}
// 	return next, index, nil
// }

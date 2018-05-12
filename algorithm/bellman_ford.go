/*
	有向图 求最短路径  bellman-ford 算法
*/
package main

import (
	// "errors"
	"fmt"
	"math"
)

// 所有路径都是双向的
var route = map[string]int{
	"ab": -1,
	"ac": 4,
	"be": 2,
	"bc": 3,
	"bd": 2,
	"dc": 5,
	"db": 1,
	"ed": -3,
}

var dist = map[string]int{
	"a": 0,
	"b": math.MaxInt32,
	"c": math.MaxInt32,
	"d": math.MaxInt32,
	"e": math.MaxInt32,
	"f": math.MaxInt32,
}

var pred = map[string]string{
	"b": "",
	"c": "",
	"d": "",
	"e": "",
	"f": "",
}

var rest = []string{"b", "c", "d", "e", "f"}

func main() {
	bellFord("a")
	fmt.Println(dist, pred)
}

func bellFord(current string) {
	for i := 0; i < len(dist)-1; i++ {
		for key, _ := range route {
			relax(key)
		}
	}
}

func relax(uv string) {
	start := string(uv[0])
	end := string(uv[1])
	weight, ok := route[uv]
	if !ok {
		return
	}
	if dist[start]+weight < dist[end] {
		dist[end] = dist[start] + weight
		pred[end] = start
	}
}

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

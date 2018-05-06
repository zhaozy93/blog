/*
	有向图 求最短路径
	s --> t  6
	s --> y  4
	y --> t  1
	t --> y  2
	t --> x  3
	y --> x  9
	y --> z  3
	z --> x  5
	x --> z  4
	z --> s  7
*/
package main

import (
	"fmt"
	"math"
)

var route = map[string]int{
	"st": 6, //s --> t  6
	"sy": 4, //s --> y  4
	"yt": 1, //y --> t  1
	"ty": 2, //t --> y  2
	"tx": 3, //t --> x  3
	"yx": 9, //y --> x  9
	"yz": 3, //y --> z  3
	"zx": 5, //z --> x  5
	"xz": 4, //x --> z  4
	"zs": 7, //z --> s  7
}

var res = map[string]int{
	"s": 0,
	"t": math.MaxInt32,
	"y": math.MaxInt32,
	"x": math.MaxInt32,
	"z": math.MaxInt32,
}

func main() {
	dijkstra("s", "x", []string{"s"}, 0)
	fmt.Println(res)
}

func dijkstra(s, e string, line []string, cost int) {
	if s == e {
		fmt.Println("找到路径", line)
		return
	}
	nextV := findNextPoint(s)
	for _, next := range nextV {
		if !isInLine(next, line) {
			w := findRouteWeight(s, next)
			// fmt.Println(s, e, w)
			if cost+w < res[next] {
				res[next] = cost + w
				fmt.Println(res)
			}

			newLine := append(line, next)
			dijkstra(next, e, newLine, res[e])
		}
	}
}

func findNextPoint(v string) []string {
	next := make([]string, 0, 0)
	for key, _ := range route {
		if key[0] == ([]byte(v))[0] {
			next = append(next, string(key[1]))
		}
	}
	return next
}

func findRouteWeight(u, v string) int {
	line := u + v
	if weight, ok := route[line]; ok {
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

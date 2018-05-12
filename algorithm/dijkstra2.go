/*
	有向图(非负) 求最短路径  dijkstra 算法
*/
package main

import (
	"errors"
	"fmt"
	"math"
)

// 所有路径都是双向的
var route = map[string]int{
	"ab": 6,
	"ac": 3,
	"bc": 2,
	"bd": 5,
	"cd": 3,
	"ce": 4,
	"de": 2,
	"df": 3,
	"ef": 5,
}

var weight = map[string]int{
	"a": 0,
	"b": math.MaxInt32,
	"c": math.MaxInt32,
	"d": math.MaxInt32,
	"e": math.MaxInt32,
	"f": math.MaxInt32,
}

var path = map[string][]string{
	"a": []string{"a"},
	"b": []string{"a"},
	"c": []string{"a"},
	"d": []string{"a"},
	"e": []string{"a"},
	"f": []string{"a"},
}

var rest = []string{"b", "c", "d", "e", "f"}

func main() {
	dijkstra("a")
	fmt.Println(weight, path)
}

func dijkstra(current string) {
	nextV := findConnPoint(current)
	line := path[current]
	cost := weight[current]
	for _, next := range nextV {
		if !isInLine(next, line) {
			w := findRouteWeight(current, next)
			if cost+w < weight[next] {
				weight[next] = cost + w
				path[next] = append(line, next)
			}
		}
	}
	if next, index, err := findNextPoint(); err == nil {
		rest = append(rest[:index], rest[index+1:]...)
		dijkstra(next)
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
	if weight, ok := route[v+u]; ok {
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

func findNextPoint() (string, int, error) {
	if len(rest) == 0 {
		return "", -1, errors.New("no point in rest arr")
	}
	next := rest[0]
	index := 0
	for i, key := range rest {
		if weight[key] < weight[next] {
			next = key
			index = i
		}
	}
	return next, index, nil
}

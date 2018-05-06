/*
	有向无环图 求各种线性排序方式
	a --> b --> c  --> e  --> f  --> g  --> n
			--> d  --> e  --> f  --> g  --> n
					   e  --> k
	h --> c
	i --> j --> k  --> m  --> n
*/
package main

import (
	"fmt"
)

var inDegree = map[string]int{
	"a": 0,
	"b": 1,
	"c": 2,
	"d": 1,
	"e": 2,
	"f": 1,
	"g": 1,
	"h": 0,
	"i": 0,
	"j": 1,
	"k": 2,
	"m": 1,
	"n": 2,
}

var outDir = map[string][]string{
	"a": []string{"b"},
	"b": []string{"c", "d"},
	"c": []string{"e"},
	"d": []string{"e"},
	"e": []string{"f", "k"},
	"f": []string{"g"},
	"g": []string{"n"},
	"h": []string{"c"},
	"i": []string{"j"},
	"j": []string{"k"},
	"k": []string{"m"},
	"m": []string{"n"},
}

var inDir = map[string][]string{
	"a": []string{},
	"b": []string{"a"},
	"c": []string{"b", "h"},
	"d": []string{"b"},
	"e": []string{"c", "d"},
	"f": []string{"e"},
	"g": []string{"f"},
	"h": []string{},
	"i": []string{},
	"j": []string{"i"},
	"k": []string{"j"},
	"m": []string{"k"},
	"n": []string{"g", "m"},
}

var result = make([][]string, 10, 10)

func main() {
	topologicalSort(inDegree, []string{})
}

func topologicalSort(uv map[string]int, line []string) {
	startV := getStartV(uv)
	// fmt.Println("当前排序状况", uv, line, startV)
	if len(startV) == 0 {
		fmt.Println("找到一个路径", line)
		return
	}
	for _, val := range startV {
		newLine := append(line, val)
		newVu := copyInDegree(uv)
		for _, v := range outDir[val] {
			newVu[v] = newVu[v] - 1
		}
		delete(newVu, val)
		topologicalSort(newVu, newLine)
	}
}

func getStartV(uv map[string]int) (arr []string) {
	for i, val := range uv {
		if val == 0 {
			arr = append(arr, i)
		}
	}
	return
}

func copyInDegree(old map[string]int) map[string]int {
	var newDegree = map[string]int{}
	for key, val := range old {
		newDegree[key] = val
	}
	return newDegree
}

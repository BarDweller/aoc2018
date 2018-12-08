package main

import (
	"fmt"
	"strconv"
	"strings"
)

type node struct {
	children []node
	metadata []int
}

func processNode(data []int, index int) (result node, after int) {
	result = node{[]node{}, []int{}}
	childCount := data[index]
	metaCount := data[index+1]
	index += 2
	for child := 0; child < childCount; child++ {
		childNode, newindex := processNode(data, index)
		index = newindex
		result.children = append(result.children, childNode)
	}
	for meta := 0; meta < metaCount; meta++ {
		result.metadata = append(result.metadata, data[index])
		index++
	}
	return result, index
}

func loaddata(input string) (root node) {

	var data []int
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		for _, val := range strings.Split(line, " ") {
			if ival, err := strconv.Atoi(val); err == nil {
				data = append(data, ival)
			} else {
				panic("Non int in line " + line)
			}
		}
	}

	index := 0
	root, index = processNode(data, index)

	return root
}

func dumpTree(n node, pad string) {
	fmt.Println(pad + " (" + strings.Join(strings.Fields(fmt.Sprint(n.metadata)), ",") + ")")
	for _, c := range n.children {
		dumpTree(c, pad+"  ")
	}
}

func sumTree(n node) int {
	thisNodeTotal := 0
	for _, m := range n.metadata {
		thisNodeTotal += m
	}
	for _, c := range n.children {
		thisNodeTotal += sumTree(c)
	}
	return thisNodeTotal
}

func main() {
	root := loaddata(testdata())
	dumpTree(root, "")
	fmt.Println(sumTree(root))
	main2(root)
}

func testdata() string {
	return `
	2 3 0 3 10 11 12 1 1 0 1 99 2 1 1 2
	`
}

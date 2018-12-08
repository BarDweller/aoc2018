package main

import "fmt"

func sumTree2(n node) int {
	thisNodeTotal := 0
	childCount := len(n.children)
	if childCount == 0 {
		for _, m := range n.metadata {
			thisNodeTotal += m
		}
	} else {
		branchCache := map[int]int{}
		for _, m := range n.metadata {
			if m <= len(n.children) {
				if branchCache[m] != 0 {
					thisNodeTotal += branchCache[m]
				} else {
					branchTotal := sumTree2(n.children[m-1])
					branchCache[m] = branchTotal
					thisNodeTotal += branchTotal
				}
			}
		}
	}
	return thisNodeTotal
}

func main2(root node) {
	fmt.Println(sumTree2(root))
}

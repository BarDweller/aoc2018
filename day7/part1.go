package main

import (
	"fmt"
	"sort"
	"strings"
)

type node struct {
	label    string
	children []*node
	parents  []*node
}

func loaddata(input string) (root node, nodemap map[string]*node) {
	nodemap = map[string]*node{}

	//build node relationships
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parent := ""
		child := ""
		fmt.Sscanf(line, "Step %s must be finished before step %s can being.", &parent, &child)
		//create child node if unknown
		if nodemap[child] == nil {
			nodemap[child] = &node{label: child, children: []*node{}, parents: []*node{}}
		}
		//create parent node if unknown
		if nodemap[parent] == nil {
			nodemap[parent] = &node{label: parent, children: []*node{}, parents: []*node{}}
		}
		//add relationships
		(*nodemap[parent]).children = append((*nodemap[parent]).children, nodemap[child])
		(*nodemap[child]).parents = append((*nodemap[child]).parents, nodemap[parent])
	}

	//apparently there can be MULTIPLE roots.. heh.. didn't see that coming!
	roots := []*node{}
	for _, node := range nodemap {
		if len(node.parents) == 0 {
			roots = append(roots, node)
		}
	}

	//in case we have multiple.. lets add our own root of "" that sits above them all
	root = node{label: "", parents: []*node{}, children: roots}
	for _, child := range roots {
		(*child).parents = append((*child).parents, &root)
	}

	return root, nodemap
}

func dumpNodes(n node, pad string) {
	next := pad + n.label + "("
	for idx, parent := range n.parents {
		if idx != 0 {
			next += ","
		}
		next += parent.label
	}
	next += ")->"
	if len(n.children) == 0 {
		fmt.Println(next)
	} else {
		for _, child := range n.children {
			dumpNodes(*child, next)
		}
	}
}

const unseen int = 0
const available int = 1
const processed int = 2

func walkNodes(chosen node, nodemap *map[string]int, nodes map[string]*node) {
	availablenodes := []string{}
	for done := false; !done; done = len(availablenodes) == 0 {
		fmt.Print(chosen.label)
		(*nodemap)[chosen.label] = processed
		for _, child := range chosen.children {
			if (*nodemap)[child.label] == unseen {
				(*nodemap)[child.label] = available
			}
		}

		availablenodes = []string{}
		for label, seen := range *nodemap {
			if seen == available {
				availablenodes = append(availablenodes, label)
			}
		}

		if len(availablenodes) > 0 {
			sort.Strings(availablenodes)

			nextNode := ""
			for _, available := range availablenodes {
				allParentsFinished := true
				for _, parent := range nodes[available].parents {
					if (*nodemap)[parent.label] != processed {
						allParentsFinished = false
						break
					}
				}
				if allParentsFinished {
					nextNode = available
					break
				}
			}
			chosen = *nodes[nextNode]
		}
	}
	fmt.Println("")
}

func main() {
	root, nodemap := loaddata(testdata())
	dumpNodes(root, "")
	visitmap := &(map[string]int{})
	walkNodes(root, visitmap, nodemap)
	main2(root, nodemap)
}

func testdata() string {
	return `
	Step C must be finished before step A can begin.
	Step C must be finished before step F can begin.
	Step A must be finished before step B can begin.
	Step A must be finished before step D can begin.
	Step B must be finished before step E can begin.
	Step D must be finished before step E can begin.
	Step F must be finished before step E can begin.	
	`
}

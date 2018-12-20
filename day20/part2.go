package main

import (
	"fmt"
	"image"

	"gopkg.in/karalabe/cookiejar.v1/collections/deque"
)

func find1000doors(linked linkedmap) int {
	start := image.Point{0, 0}
	type node struct {
		distance int
		coord    image.Point
	}
	initial := node{0, start}
	q := deque.New()
	q.PushRight(initial)
	count := 0
	seen := map[image.Point]bool{start: true}
	for !q.Empty() {
		n := q.PopLeft().(node)
		if n.distance >= 1000 {
			count++
		}
		for delta := range linked[n.coord] {
			dn := n.coord.Add(delta)
			if seen[dn] {
				continue
			}
			seen[dn] = true
			q.PushRight(node{n.distance + 1, dn})
		}
	}
	return count
}

func part2(input string) {
	lm := buildlinkedmap(input)
	fmt.Println("1000 doors count", find1000doors(lm))
}

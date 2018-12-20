package main

import (
	"fmt"
	"image"

	"gopkg.in/karalabe/cookiejar.v1/collections/deque"
)

type linkedmap map[image.Point]map[image.Point]bool

func buildlinkedmap(input string) linkedmap {
	start := image.Point{0, 0}
	s := deque.New()
	s.PushRight(start)
	linked := linkedmap{}
	cur := image.Point{0, 0}
	dmap := map[rune]image.Point{'N': image.Point{0, -1}, 'S': image.Point{0, 1}, 'E': image.Point{1, 0}, 'W': image.Point{-1, 0}}
	for _, r := range input {
		switch r {
		case 'N', 'S', 'E', 'W':
			{
				delta := dmap[r]
				if _, found := linked[cur]; !found {
					linked[cur] = map[image.Point]bool{}
				}
				linked[cur][delta] = true
				cur = cur.Add(delta)
				backwards := delta.Mul(-1)
				if _, found := linked[cur]; !found {
					linked[cur] = map[image.Point]bool{}
				}
				linked[cur][backwards] = true
			}
		case '(':
			s.PushRight(cur)
		case ')':
			s.PopRight()
		case '|':
			cur = s.Right().(image.Point)
		}
	}
	return linked
}

func findmaxdoors(linked linkedmap) int {
	start := image.Point{0, 0}
	maxdistance := 0
	type node struct {
		distance int
		coord    image.Point
	}
	initial := node{0, start}
	q := deque.New()
	q.PushRight(initial)
	seen := map[image.Point]bool{start: true}
	for !q.Empty() {
		n := q.PopLeft().(node)
		if n.distance > maxdistance {
			maxdistance = n.distance
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
	return maxdistance
}

func part1(input string) {
	lm := buildlinkedmap(input)
	fmt.Println("Max doors ", findmaxdoors(lm))
}

func main() {
	part1(testdata1())
	part2(testdata2())
}

func testdata1() string {
	return `^WSSEESWWWNW(S|NENNEEEENN(ESSSSW(NWSW|SSEN)|WSWWN(E|WWS(E|SS))))$`
}

func testdata2() string {
	return `^ESSWWN(E|NNENN(EESS(WNSE|)SSS|WWWSSSSE(SW|NNNE)))$`
}

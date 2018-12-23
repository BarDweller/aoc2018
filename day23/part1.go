package main

import (
	"fmt"
	"strings"
)

type coord struct {
	x, y, z int
}

type nanobot struct {
	loc    coord
	radius int
}

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

func manhattandistance(a coord, b coord) int {
	return abs(a.x-b.x) + abs(a.y-b.y) + abs(a.z-b.z)
}

func loaddata(input string) []nanobot {
	result := []nanobot{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n := nanobot{}
		fmt.Sscanf(line, "pos=<%d,%d,%d>, r=%d", &n.loc.x, &n.loc.y, &n.loc.z, &n.radius)
		result = append(result, n)
	}
	return result
}

func reachablecount(chosennano nanobot, nanos []nanobot) int {
	count := 0
	for _, nano := range nanos {
		if manhattandistance(nano.loc, chosennano.loc) <= chosennano.radius {
			count++
		}
	}
	return count
}

func part1(input string) {
	nanos := loaddata(input)
	maxr := 0
	var chosennano nanobot
	for _, nano := range nanos {
		if nano.radius > maxr {
			maxr = nano.radius
			chosennano = nano
		}
	}
	fmt.Println("selected nano ", chosennano)
	count := reachablecount(chosennano, nanos)
	fmt.Println("There are", count, "nanobots in range")
}

func main() {
	part1(data())
	part2(data())
}

func testdata() string {
	return `
	pos=<0,0,0>, r=4
	pos=<1,0,0>, r=1
	pos=<4,0,0>, r=3
	pos=<0,2,0>, r=1
	pos=<0,5,0>, r=3
	pos=<0,0,3>, r=1
	pos=<1,1,1>, r=1
	pos=<1,1,2>, r=1
	pos=<1,3,1>, r=1	
	`
}

func testdata2() string {
	return `
	pos=<10,12,12>, r=2
	pos=<12,14,12>, r=2
	pos=<16,12,12>, r=4
	pos=<14,14,14>, r=6
	pos=<50,50,50>, r=200
	pos=<10,10,10>, r=5	
	`
}

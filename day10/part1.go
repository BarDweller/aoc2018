package main

import (
	"fmt"
	"math"
	"strings"
)

type particle struct {
	x, y, dx, dy int
}

func loaddata(input string) (result []particle) {
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		p := particle{}
		fmt.Sscanf(line, "position=<%d, %d> velocity=<%d, %d>", &p.x, &p.y, &p.dx, &p.dy)
		result = append(result, p)
	}
	return result
}

func moveparticles(ps []particle) {
	for i, p := range ps {
		p.x += p.dx
		p.y += p.dy
		ps[i] = p
	}
}

type xy struct {
	x, y int
}
type boundedgrid struct {
	minx, miny, maxx, maxy int
	grid                   map[xy]bool
}

func buildgrid(particles []particle) boundedgrid {
	grid := map[xy]bool{}
	minx := math.MaxInt32
	miny := math.MaxInt32
	maxx := 0
	maxy := 0
	for _, p := range particles {
		grid[xy{x: p.x, y: p.y}] = true
		if p.x > maxx {
			maxx = p.x
		}
		if p.y > maxy {
			maxy = p.y
		}
		if p.x < minx {
			minx = p.x
		}
		if p.y < miny {
			miny = p.y
		}
	}
	return boundedgrid{minx: minx, miny: miny, maxx: maxx, maxy: maxy, grid: grid}
}

//dump grid using same output as example
func printgrid(particles []particle) {
	bgrid := buildgrid(particles)
	for y := bgrid.miny; y <= bgrid.maxy; y++ {
		for x := bgrid.minx; x <= bgrid.maxx; x++ {
			if bgrid.grid[xy{x: x, y: y}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
}

func candidate(particles []particle) bool {
	//look for at least 3 sets of 6 #'s aligned vertically
	seekcols := 3
	bgrid := buildgrid(particles)
	cols := 0
	for t := range bgrid.grid {
		if bgrid.grid[xy{t.x, t.y - 1}] &&
			bgrid.grid[xy{t.x, t.y - 2}] &&
			bgrid.grid[xy{t.x, t.y - 3}] &&
			bgrid.grid[xy{t.x, t.y - 4}] &&
			bgrid.grid[xy{t.x, t.y - 5}] {
			cols++
		}
		if cols > seekcols-1 {
			break
		}
	}
	return cols > seekcols-1
}

func moveuntilcandidate(particles []particle) int {
	iter := 0
	found := false
	for !found {
		iter++
		moveparticles(particles)
		found = candidate(particles)
	}
	printgrid(particles)
	return iter
}

func main() {
	particles := loaddata(testdata())
	iter := moveuntilcandidate(particles)
	main2(iter)
}

func testdata() string {
	return `
	position=< 9,  1> velocity=< 0,  2>
	position=< 7,  0> velocity=<-1,  0>
	position=< 3, -2> velocity=<-1,  1>
	position=< 6, 10> velocity=<-2, -1>
	position=< 2, -4> velocity=< 2,  2>
	position=<-6, 10> velocity=< 2, -2>
	position=< 1,  8> velocity=< 1, -1>
	position=< 1,  7> velocity=< 1,  0>
	position=<-3, 11> velocity=< 1, -2>
	position=< 7,  6> velocity=<-1, -1>
	position=<-2,  3> velocity=< 1,  0>
	position=<-4,  3> velocity=< 2,  0>
	position=<10, -3> velocity=<-1,  1>
	position=< 5, 11> velocity=< 1, -2>
	position=< 4,  7> velocity=< 0, -1>
	position=< 8, -2> velocity=< 0,  1>
	position=<15,  0> velocity=<-2,  0>
	position=< 1,  6> velocity=< 1,  0>
	position=< 8,  9> velocity=< 0, -1>
	position=< 3,  3> velocity=<-1,  1>
	position=< 0,  5> velocity=< 0, -1>
	position=<-2,  2> velocity=< 2,  0>
	position=< 5, -2> velocity=< 1,  2>
	position=< 1,  4> velocity=< 2,  1>
	position=<-2,  7> velocity=< 2, -2>
	position=< 3,  6> velocity=<-1, -1>
	position=< 5,  0> velocity=< 1,  0>
	position=<-6,  0> velocity=< 2,  0>
	position=< 5,  9> velocity=< 1, -2>
	position=<14,  7> velocity=<-2,  0>
	position=<-3,  6> velocity=< 2, -1>	
	`
}

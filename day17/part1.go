package main

import (
	"fmt"
	"strings"
)

type xy struct {
	x, y int
}

const (
	springx = 500
	springy = 0
)

type state string

const (
	sand         = state(".")
	clay         = state("#")
	waterstopped = state("~")
	watermotion  = state("|")
	spring       = state("+")
)

type grid struct {
	d   map[xy]state
	max xy
	min xy
}

func loaddata(input string) grid {
	g := grid{map[xy]state{}, xy{0, 0}, xy{0, 0}}
	//add fixed spring
	g.d[xy{springx, springy}] = spring
	maxx, maxy := 0, 0
	minx, miny := 4096, 4096
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line[0] == 'x' {
			x, y1, y2 := 0, 0, 0
			fmt.Sscanf(line, "x=%d, y=%d..%d", &x, &y1, &y2)
			for y := y1; y <= y2; y++ {
				g.d[xy{x, y}] = clay
			}
			if x > maxx {
				maxx = x
			}
			if y2 > maxy {
				maxy = y2
			}
			if x < minx {
				minx = x
			}
			if y1 < miny {
				miny = y1
			}
		}
		if line[0] == 'y' {
			y, x1, x2 := 0, 0, 0
			fmt.Sscanf(line, "y=%d, x=%d..%d", &y, &x1, &x2)
			for x := x1; x <= x2; x++ {
				g.d[xy{x, y}] = clay
			}
			if x2 > maxx {
				maxx = x2
			}
			if y > maxy {
				maxy = y
			}
			if x1 < minx {
				minx = x1
			}
			if y < miny {
				miny = y
			}
		}
	}
	g.max = xy{maxx, maxy}
	g.min = xy{minx, miny}
	return g
}

func followfall(w xy, g grid, widths map[xy]bool) {
	//last := w
	for w.y < g.max.y {
		//ppgrid(g)
		//look beneath current sq.
		below := xy{w.x, w.y + 1}
		switch g.d[below] {
		case waterstopped, clay:
			//only process falls we haven't seen before.
			if _, found := widths[below]; found {
				//we've already done this waterfall.
				return
			} else {
				//mark that we've processed this waterfall.
				widths[below] = true

				//create width for this fall
				leftcount := 0
				leftwall := false
				rightcount := 0
				rightwall := false

				//current cell is in motion..(unless already set to something)
				if _, found := g.d[w]; !found {
					g.d[w] = watermotion
				}

				//check for boundaries..
				//a) count left..
				for x := w.x - 1; x >= g.min.x-1; x-- {
					//needs to be bounded below to continue loop
					scanbelow := xy{x, below.y}
					if !(g.d[scanbelow] == clay || g.d[scanbelow] == waterstopped) {
						if _, found := g.d[xy{x, w.y}]; !found {
							g.d[xy{x, w.y}] = watermotion
						}
						followfall(xy{x, w.y}, g, widths)
						break
					}
					//needs x not to be clay to continue loop
					rowx := xy{x, w.y}
					if g.d[rowx] == clay {
						leftwall = true
						break
					}
					//set as in motion (unless already set to something)
					if _, found := g.d[rowx]; !found {
						g.d[rowx] = watermotion
					}
					leftcount++
				}
				//b) count right..
				for x := w.x + 1; x <= g.max.x+1; x++ {
					//needs to be bounded below to continue loop
					scanbelow := xy{x, below.y}
					if !(g.d[scanbelow] == clay || g.d[scanbelow] == waterstopped) {
						if _, found := g.d[xy{x, w.y}]; !found {
							g.d[xy{x, w.y}] = watermotion
						}
						followfall(xy{x, w.y}, g, widths)
						break
					}
					//needs x not to be clay to continue loop
					rowx := xy{x, w.y}
					if g.d[rowx] == clay {
						rightwall = true
						break
					}
					//set as in motion (unless already set to something)
					if _, found := g.d[rowx]; !found {
						g.d[rowx] = watermotion
					}
					rightcount++
				}

				//if the waterfall led to a bounded space, convert it to
				//stationary water.
				if leftwall && rightwall {
					for x := w.x - leftcount; x < w.x+rightcount+1; x++ {
						//stationary water overwrites water in motion,
						//so no need to test.
						g.d[xy{x, w.y}] = waterstopped
					}
					//we filled in a row, backtrack up a row & re-evaluate
					w = xy{w.x, w.y - 1}
				} else {
					//if we had a missing wall left or right,
					//then we already followed the falls as required
					//so we are done following 'this' fall.
					return
				}
			}
		case watermotion:
			//no op, we'll keep falling down.
			w = below
		default: //sand.
			g.d[below] = watermotion
			w = below
		}
	}
}

func count(g grid) int {
	total := 0
	for y := g.min.y; y <= g.max.y+1; y++ {
		for x := g.min.x - 1; x <= g.max.x+1; x++ {
			if state, present := g.d[xy{x, y}]; present {
				switch state {
				case watermotion, waterstopped:
					total++
				}
			}
		}
	}
	return total
}

func ppgrid(g grid) {
	fmt.Println(g.min, g.max)
	for y := g.min.y - 1; y <= g.max.y+1; y++ {
		line := ""
		for x := g.min.x - 1; x <= g.max.x+1; x++ {
			if state, present := g.d[xy{x, y}]; present {
				line += string(state)
			} else {
				line += "."
			}
		}
		fmt.Println(line)
	}
	fmt.Println("")
}

func part1(g grid) {
	widths := map[xy]bool{}
	w := xy{springx, springy}
	followfall(w, g, widths)
	ppgrid(g)
	watercells := count(g)
	fmt.Println("Water: ", watercells)
}

func main() {
	data := loaddata(testdata())
	part1(data)
	main2(data)
}

func testdata() string {
	return `
x=495, y=2..7
y=7, x=495..501
x=501, y=3..7
x=498, y=2..4
x=506, y=1..2
x=498, y=10..13
x=504, y=10..13
y=13, x=498..504
`
}

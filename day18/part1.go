package main

import (
	"fmt"
	"strings"
)

type xy struct {
	x, y int
}

type grid struct {
	d   map[xy]rune
	max xy
}

func loaddata(input string) grid {
	g := grid{map[xy]rune{}, xy{0, 0}}
	y := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		x := 0
		for _, r := range line {
			g.d[xy{x, y}] = r
			x++
			if x > g.max.x {
				g.max.x = x
			}
		}
		y++
		if y > g.max.y {
			g.max.y = y
		}
	}
	return g
}

func pp(g grid) {
	for y := 0; y < g.max.y; y++ {
		for x := 0; x < g.max.x; x++ {
			fmt.Print(string(g.d[xy{x, y}]))
		}
		fmt.Println("")
	}
	fmt.Println("")
}

type mutationresult struct {
	cell   xy
	result rune
}

func mutatecell(cell xy, g grid, r chan mutationresult) {
	neighbours := []xy{xy{cell.x - 1, cell.y},
		xy{cell.x - 1, cell.y - 1},
		xy{cell.x, cell.y - 1},
		xy{cell.x + 1, cell.y - 1},
		xy{cell.x + 1, cell.y},
		xy{cell.x + 1, cell.y + 1},
		xy{cell.x, cell.y + 1},
		xy{cell.x - 1, cell.y + 1}}

	result := g.d[cell]
	switch g.d[cell] {
	case '.', '|':
		{
			seek := '|'
			if g.d[cell] == '|' {
				seek = '#'
			}
			total := 0
			for _, c := range neighbours {
				if n, found := g.d[c]; found && n == seek {
					total++
					if total >= 3 {
						if result == '.' {
							result = '|'
						} else {
							result = '#'
						}
						break
					}
				}
			}
		}
	case '#':
		{
			tree := false
			mill := false
			for _, c := range neighbours {
				if n, found := g.d[c]; found && (n == '#' || n == '|') {
					if !tree {
						tree = n == '|'
					}
					if !mill {
						mill = n == '#'
					}
				}
			}
			if !(mill && tree) {
				result = '.'
			}
		}
	}
	r <- mutationresult{cell, result}
}

func doMinute(g grid) grid {
	results := make(chan mutationresult)
	for k := range g.d {
		go mutatecell(k, g, results)
	}
	ng := map[xy]rune{}
	for count := len(g.d); count > 0; count-- {
		mutationresult := <-results
		ng[mutationresult.cell] = mutationresult.result
	}
	g.d = ng
	return g
}

func count(g grid) {
	c := map[rune]int{}
	for _, v := range g.d {
		c[v]++
	}
	fmt.Println("Open: ", c['.'])
	fmt.Println("Wood: ", c['|'])
	fmt.Println("Mill: ", c['#'])

	fmt.Println("Resource Value: ", c['#']*c['|'])
}

func part1(g grid) {
	for i := 1; i <= 10; i++ {
		g = doMinute(g)
    //enable for pretty print as per example ;)
		//fmt.Println("After", i, "minutes:")
		//pp(g)
	}
	count(g)
}

func main() {
	g := loaddata(testdata())
	part1(g)

	main2()
}

func testdata() string {
	return `
	.#.#...|#.
	.....#|##|
	.|..|...#.
	..|#.....#
	#.#|||#|#|
	...#.||...
	.|....|...
	||...#|.#|
	|.||||..|.
	...#.|..|.	
	`
}

package main

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/karalabe/cookiejar.v1/collections/deque"
)

type race byte

const (
	elf    race = 0
	goblin race = 1
)

const (
	default_hp          int = 200
	default_attackpower int = 3
)

type unit struct {
	side        race
	hitpoints   int
	attackpower int
	xy          coord
}

func (u unit) String() string {
	side := "G"
	if u.side == elf {
		side = "E"
	}
	return fmt.Sprintf("{%s:hp:%d:%s}", side, u.hitpoints, u.xy.String())
}

type coord struct {
	x, y int
}

func (c coord) String() string {
	return fmt.Sprintf("{%d,%d}", c.x, c.y)
}

type grid map[coord]rune

func readingOrder(a, b unit) unit {
	switch {
	case a.xy.y < b.xy.y:
		return a
	case a.xy.y > b.xy.y:
		return b
	default:
		if a.xy.x < b.xy.x {
			return a
		} else {
			return b
		}
	}
}

func sortByReadingOrder(units []unit) {
	sort.Slice(units, func(i, j int) bool {
		switch {
		case units[i].xy.y < units[j].xy.y:
			return true
		case units[i].xy.y > units[j].xy.y:
			return false
		default:
			return units[i].xy.x < units[j].xy.x
		}
	})
}

func sortCoordsByReadingOrder(coords []coord) {
	sort.Slice(coords, func(i, j int) bool {
		switch {
		case coords[i].y < coords[j].y:
			return true
		case coords[i].y > coords[j].y:
			return false
		default:
			return coords[i].x < coords[j].x
		}
	})
}

func adjacent(a, b unit) bool {
	return (a.xy.x == b.xy.x && (a.xy.y == b.xy.y-1 || a.xy.y == b.xy.y+1)) ||
		(a.xy.y == b.xy.y && (a.xy.x == b.xy.x-1 || a.xy.x == b.xy.x+1))
}

func selectWeakestEnemy(enemies []unit) unit {
	selected := enemies[0]
	if len(enemies) == 1 {
		return selected
	}
	for _, e := range enemies[1:] {
		if e.hitpoints == selected.hitpoints {
			selected = readingOrder(selected, e)
		}
		if e.hitpoints < selected.hitpoints {
			selected = e
		}
	}
	return selected
}

func clearNeighbours(p coord, g grid) []coord {
	clear := []coord{}
	up := coord{p.x, p.y - 1}
	down := coord{p.x, p.y + 1}
	left := coord{p.x - 1, p.y}
	right := coord{p.x + 1, p.y}
	if g[up] == '.' {
		clear = append(clear, up)
	}
	if g[left] == '.' {
		clear = append(clear, left)
	}
	if g[right] == '.' {
		clear = append(clear, right)
	}
	if g[down] == '.' {
		clear = append(clear, down)
	}
	return clear
}

type node struct {
	xy       coord
	distance int
}

func (n node) String() string {
	return fmt.Sprintf("%s:%d", n.xy.String(), n.distance)
}

func findClosest(g grid, start coord, targets []coord) ([]coord, int) {
	tmap := map[coord]bool{}
	for _, t := range targets {
		tmap[t] = true
	}
	result := []coord{}
	q := deque.New()
	q.PushRight(node{start, 0})
	distance := -1
	seen := map[coord]bool{}
	for !q.Empty() {
		current := q.PopLeft().(node)
		if distance != -1 && current.distance > distance {
			return result, distance
		}
		if seen[current.xy] {
			continue
		}
		seen[current.xy] = true
		if tmap[current.xy] {
			distance = current.distance
			result = append(result, current.xy)
		}
		for _, c := range clearNeighbours(current.xy, g) {
			if !seen[c] {
				q.PushRight(node{c, current.distance + 1})
			}
		}
	}
	return result, distance
}

func round(units []unit, g grid, x int, y int) bool {
	//sort units by reading order (y, then x)
	sortByReadingOrder(units)

	//for each unit
	for i := range units {
		u := units[i]
		if u.hitpoints <= 0 {
			continue
		}

		if warover(units) {
			return true
		}

		//identify all targets (enemies)
		enemies := []unit{}
		for _, pe := range units {
			//ignore dead enemies =)
			if pe.side != u.side && pe.hitpoints > 0 {
				enemies = append(enemies, pe)
			}
		}

		//identify if unit is adjacent(u/d/l/r) of enemy
		adjacentenemies := []unit{}
		for _, e := range enemies {
			if adjacent(u, e) {
				adjacentenemies = append(adjacentenemies, e)
			}
		}

		if len(adjacentenemies) == 0 {
			//  for each enemy unit, identify adjacent(u/d/l/r) empty (.) squares
			targetcoords := []coord{}
			for _, ae := range enemies {
				targetcoords = append(targetcoords, clearNeighbours(ae.xy, g)...)
			}
			//    if no empty squares beside enemies, end turn
			if len(targetcoords) == 0 {
				continue
			}

			//  for each empty square
			//    determine path length
			closest, distance := findClosest(g, u.xy, targetcoords)

			//  if all squares are unreachable, end turn
			if len(closest) == 0 {
				continue
			}

			sortCoordsByReadingOrder(closest)
			chosen := []coord{closest[0]}

			for _, c := range clearNeighbours(u.xy, g) {
				_, d := findClosest(g, c, chosen)
				if d == distance-1 {
					g[u.xy] = '.'
					units[i].xy = c
					if u.side == elf {
						g[c] = 'E'
					} else {
						g[c] = 'G'
					}
					break
				}
			}

			//update adjacent enemies after move.
			for _, e := range enemies {
				if adjacent(units[i], e) {
					adjacentenemies = append(adjacentenemies, e)
				}
			}
		}

		//after move phase, attack
		if len(adjacentenemies) == 0 {
			//if no adjacent enemies, end turn.
			continue
		}

		//select weakest enemy
		chosen := selectWeakestEnemy(adjacentenemies)
		for ce := range units {
			//ideally we'd do chosen.hitpoints here, but
			//chosen appears to be a copy, not the item in the array
			//so go find the matching entry in the array.
			//it MUST have the same coords, but ignore already dead
			//units because they may have same coords as live ones!
			if units[ce].xy == chosen.xy && units[ce].hitpoints > 0 {
				units[ce].hitpoints -= u.attackpower
				if units[ce].hitpoints <= 0 {
					g[chosen.xy] = '.'
				}
				break
			}
		}
	}
	return false
}

func warover(units []unit) bool {
	elves := false
	goblins := false
	for _, u := range units {
		if u.hitpoints > 0 {
			if u.side == elf {
				elves = true
			} else {
				goblins = true
			}
		}
		if elves && goblins {
			break
		}
	}
	return (elves && !goblins) || (goblins && !elves)
}

func sumhitpoints(units []unit) int {
	total := 0
	for _, u := range units {
		if u.hitpoints > 0 {
			total += u.hitpoints
		}
	}
	return total
}

func pp(g grid, mx, my int) {
	for y := 0; y <= my; y++ {
		for x := 0; x <= mx; x++ {
			fmt.Print(string(g[coord{x, y}]))
		}
		fmt.Println("")
	}
}

func part1(units []unit, g grid, x int, y int) {
	var currentround int
	for currentround = 1; !warover(units); currentround++ {
		partialround := round(units, g, x, y)
		if partialround {
			//lastround didn't count
			currentround--
		}
	}
	//round is still ++ by forloop, so -1 to go back to required value
	fmt.Println("last full round ", currentround-1, (currentround-1)*sumhitpoints(units))
}

func loaddata(input string) ([]unit, grid, int, int) {
	grid := grid{}
	units := []unit{}
	x := 0
	y := 0
	maxx := 0
	maxy := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		for _, r := range line {
			grid[coord{x, y}] = r
			if r == 'G' {
				units = append(units, unit{goblin, default_hp, default_attackpower, coord{x, y}})
			} else if r == 'E' {
				units = append(units, unit{elf, default_hp, default_attackpower, coord{x, y}})
			}
			maxx = x
			x++
		}
		maxy = y
		y++
		x = 0
	}
	return units, grid, maxx, maxy
}

func main() {
	fmt.Println(">Part 1")
	units, g, x, y := loaddata(testdata())
	part1(units, g, x, y)
	main2()
}

func testdata() string {
	return `
	#######
	#.G...#
	#...EG#
	#.#.#G#
	#..G#E#
	#.....#
	#######
	`
}

func testdata2() string {
	return `
	#######
	#G..#E#
	#E#E.E#
	#G.##.#
	#...#E#
	#...E.#
	#######
	`
}

func testdata3() string {
	return `
	#######
	#E..EG#
	#.#G.E#
	#E.##E#
	#G..#.#
	#..E#.#
	#######
	`
}

func testdata5() string {
	return `
	#######
	#.E...#
	#.#..G#
	#.###.#
	#E#G#G#
	#...#G#
	#######
	`
}

func testdata6() string {
	return `
	#########
	#G......#
	#.E.#...#
	#..##..G#
	#...##..#
	#...#...#
	#.G...G.#
	#.....G.#
	#########
	`
}


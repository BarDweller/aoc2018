package main

import (
	"fmt"
	"strings"
)

func loaddata2(input string, elfpower int) ([]unit, grid, int, int) {
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
				units = append(units, unit{elf, default_hp, elfpower, coord{x, y}})
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

func round2(units []unit, g grid, x int, y int) (bool, bool) {
	//sort units by reading order (y, then x)
	sortByReadingOrder(units)

	//for each unit
	for i := range units {
		u := units[i]
		if u.hitpoints <= 0 {
			continue
		}

		if warover(units) {
			return true, false
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
					if units[ce].side == elf {
						return false, true
					}
				}
				break
			}
		}
	}
	return false, false
}

func part2(units []unit, g grid, x int, y int) bool {
	var currentround int
	for currentround = 1; !warover(units); currentround++ {
		partialround, elfdied := round2(units, g, x, y)
		if elfdied {
			return false
		}
		if partialround {
			//this round didn't count
			currentround--
		}
	}
	//round will still be ++ due to for loop, so -1 to get last value
	fmt.Println("last full round ", currentround-1, (currentround-1)*sumhitpoints(units))
	return true

}

func main2() {
	fmt.Println(">Part 2")
	for elfpower := 4; elfpower < 200; elfpower++ {
		fmt.Println("Testing elf powah ", elfpower)
		units, g, x, y := loaddata2(data(), elfpower)
		if part2(units, g, x, y) {
			fmt.Println("Elf Power was ", elfpower)
			break
		}
	}
}

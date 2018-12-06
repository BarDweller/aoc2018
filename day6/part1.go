package main

import (
	"fmt"
	"strings"
)

type coord struct {
	x, y int
}

func loaddata(input string) []coord {
	coords := []coord{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		c := coord{}
		fmt.Sscanf(line, "%d, %d", &c.x, &c.y)
		coords = append(coords, c)
	}
	return coords
}

func discoverbounds(all []coord) (lx, hx, ly, hy int) {

	hx, hy = 0, 0
	//slight cheat, looking at data we know nothing is larger
	//although we could scan the data & figure out the largest x & y
	//to use, it's not worth it =)
	lx, ly = 1024, 1024

	for _, c := range all {
		if lx > c.x {
			lx = c.x
		}
		if ly > c.y {
			ly = c.y
		}
		if hx < c.x {
			hx = c.x
		}
		if hy < c.y {
			hy = c.y
		}
	}

	return lx, hx, ly, hy
}

type griddata struct {
	index, distance int
	equidistant     bool
}

func creategrid(lx, hx, ly, hy int, all []coord) map[coord]griddata {
	grid := map[coord]griddata{}

	//init grid to be max distance everywhere for imaginary coord idx -1
	maxdist := (hx - lx) + (hy - ly)
	for y := ly; y <= hy; y++ {
		for x := lx; x <= hx; x++ {
			grid[coord{x, y}] = griddata{-1, maxdist, false}
		}
	}

	fmt.Println("Grid size : ", hx-lx, hy-ly)

	//populate grid with coord info.
	for idx, c := range all {
		floodmap(grid, idx, c, lx, hx, ly, hy)

		//fmt.Println("After index ", idx)
		//ppGrid2(grid, lx, hx, ly, hy)
	}

	return grid
}

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

func manhattandistance(a coord, b coord) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

//return a set of coords for a given ring based around an anchor coord (center of ring)
func getcoordsforring(ringcount int, anchor coord) []coord {
	if ringcount < 1 {
		panic("fiery ball of death")
	}
	perside := (ringcount * 2) + 1
	offset := ringcount
	ring := []coord{}
	//top & bottom
	for x := -offset; x < (perside - offset); x++ {
		ring = append(ring, coord{anchor.x + x, anchor.y + offset})
		ring = append(ring, coord{anchor.x + x, anchor.y - offset})
	}
	//sides
	for y := -offset + 1; y < (perside - 2 - offset + 1); y++ {
		ring = append(ring, coord{anchor.x + offset, anchor.y + y})
		ring = append(ring, coord{anchor.x - offset, anchor.y + y})
	}
	return ring
}

func inside(test coord, lx, hx, ly, hy int) bool {
	return test.x >= lx && test.x <= hx && test.y >= ly && test.y <= hy
}

//place the index 'idx' into the grid at location 'c' and ripple outward it's distances
//continuing as long as the distances are in the grid, and smaller than existing distances
func floodmap(grid map[coord]griddata, idx int, c coord, lx, hx, ly, hy int) {
	ring := 0
	insidegrid := true  //bugout when ring is no longer contained in grid bounds
	updatedgrid := true //bugout if ring resulted in no updates (existing nodes are closer)

	//ring 0 is mandatory.. it's the node itself.
	grid[c] = griddata{idx, 0, false}

	//for the rest of the rings
	for ring = 1; insidegrid && updatedgrid; ring++ {
		ringcoords := getcoordsforring(ring, c)
		insidecount := len(ringcoords)
		updatedcount := 0
		for _, rc := range ringcoords {
			if !inside(rc, lx, hx, ly, hy) {
				insidecount--
				continue
			}
			distance := manhattandistance(c, rc)
			if distance < grid[rc].distance {
				grid[rc] = griddata{idx, distance, false}
				updatedcount++
			} else if distance == grid[rc].distance {
				grid[rc] = griddata{grid[rc].index, grid[rc].distance, true}
				updatedcount++
			}
		}
		if insidecount == 0 {
			insidegrid = false
		}
		if updatedcount == 0 {
			updatedgrid = false
		}
	}

}

//build map of index->allocated cells, using rules from day6
func countallocations(grid map[coord]griddata, all []coord, lx, hx, ly, hy int) map[int]int {
	//identify edge nodes.
	interior := map[int]bool{}
	for idx, c := range all {
		if c.x == lx || c.x == hx || c.y == ly || c.y == hy {
			continue
		}
		interior[idx] = true
	}

	allocations := map[int]int{}
	//iterate grid and aggregate count for each coord idx
	//skipping edge nodes (because their counts are inaccurate, as they are infinite really)
	//skipping equidistant nodes
	for _, data := range grid {
		if interior[data.index] && !data.equidistant {
			allocations[data.index]++
		}
	}

	return allocations
}

//find the max value in the map
func findmax(allocations map[int]int) int {
	max := 0
	for _, a := range allocations {
		if a > max {
			max = a
		}
	}
	return max
}

//print the letters with their distances (for debug)
func ppGrid(grid map[coord]griddata, lx, hx, ly, hy int) {
	for y := ly; y <= hy; y++ {
		for x := lx; x <= hx; x++ {
			data := grid[coord{x, y}]
			letter := 'a' + data.index
			switch {
			case data.index == -1:
				{
					fmt.Printf(" x%03dx ", data.distance)
				}
			case data.distance == 0:
				{
					fmt.Printf("<%c%03d%c>", letter, data.index, letter)
				}
			case data.equidistant:
				{
					fmt.Printf(" .%03d. ", data.distance)
				}
			default:
				{
					fmt.Printf(" %c%03d%c ", letter, data.distance, letter)
				}
			}
		}
		fmt.Println("")
	}
}

//print the letters like the example does
//assumes we have less than 26 indexes.. if not, we start
//using extended ascii ;) looks pretty tho!
func ppGrid2(grid map[coord]griddata, lx, hx, ly, hy int) {
	for y := ly; y <= hy; y++ {
		for x := lx; x <= hx; x++ {
			data := grid[coord{x, y}]
			letter := 'a' + data.index
			ucase := 'A' + data.index
			switch {
			case data.index == -1:
				{
					fmt.Printf("*")
				}
			case data.distance == 0:
				{
					fmt.Printf("%c", ucase)
				}
			case data.equidistant:
				{
					fmt.Printf(".")
				}
			default:
				{
					fmt.Printf("%c", letter)
				}
			}
		}
		fmt.Println("")
	}
}

func main() {
	coords := loaddata(testdata())
	lx, hx, ly, hy := discoverbounds(coords)
	grid := creategrid(lx, hx, ly, hy, coords)
	ppGrid2(grid, lx, hx, ly, hy)
	allocations := countallocations(grid, coords, lx, hx, ly, hy)
	max := findmax(allocations)
	fmt.Println(max)
	main2(coords, lx, hx, ly, hy)
}

func testdata() string {
	return `
	1, 1
	1, 6
	8, 3
	3, 4
	5, 5
	8, 9	
	`
}

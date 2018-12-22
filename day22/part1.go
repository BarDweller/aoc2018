package main

import "fmt"

type xy struct {
	x, y int
}

type grid struct {
	target       xy
	depth        int
	geologicbyxy map[xy]int
	erosionbyxy  map[xy]int
}

func getgeologic(t xy, g grid) int {
	zero := xy{0, 0}
	if geo, found := g.geologicbyxy[t]; found {
		return geo
	}
	switch {
	case t == zero, t == g.target:
		g.geologicbyxy[t] = 0
	case t.x == 0:
		g.geologicbyxy[t] = t.y * 48271
	case t.y == 0:
		g.geologicbyxy[t] = t.x * 16807
	default:
		g.geologicbyxy[t] = geterosion(xy{t.x - 1, t.y}, g) * geterosion(xy{t.x, t.y - 1}, g)
	}
	return g.geologicbyxy[t]
}

func geterosion(t xy, g grid) int {
	if ero, found := g.erosionbyxy[t]; found {
		return ero
	}
	g.erosionbyxy[t] = (getgeologic(t, g) + g.depth) % 20183
	return g.erosionbyxy[t]
}

func pp(g grid, mx, my int) {
	for y := 0; y < my; y++ {
		for x := 0; x < mx; x++ {
			erosion := geterosion(xy{x, y}, g)
			switch {
			case x == 0 && y == 0:
				fmt.Print("M")
			case xy{x, y} == g.target:
				fmt.Print("T")
			default:
				switch erosion % 3 {
				case 0:
					fmt.Print(".")
				case 1:
					fmt.Print("=")
				case 2:
					fmt.Print("|")
				}
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func risklevel(g grid) int {
	sum := 0
	for y := 0; y <= g.target.y; y++ {
		for x := 0; x <= g.target.x; x++ {
			sum += geterosion(xy{x, y}, g) % 3
		}
	}
	return sum
}

func main() {
	//example
	g := grid{xy{10, 10}, 510, map[xy]int{}, map[xy]int{}}

	//actual input
	//g := grid{xy{13, 734}, 7305, map[xy]int{}, map[xy]int{}}

	//pp(g, 16, 16)
	fmt.Println("Risk", risklevel(g))

	part2(g)
}

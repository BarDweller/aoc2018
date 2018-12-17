package main

import "fmt"

func count2(g grid) int {
	total := 0
	for y := g.min.y; y <= g.max.y+1; y++ {
		for x := g.min.x - 1; x <= g.max.x+1; x++ {
			if state, present := g.d[xy{x, y}]; present {
				switch state {
				case waterstopped:
					total++
				}
			}
		}
	}
	return total
}

func main2(g grid) {
	fmt.Println("Stationary Water : ", count2(g))
}

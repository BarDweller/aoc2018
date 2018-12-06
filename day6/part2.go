package main

import (
	"fmt"
)

func calcdistance(c coord, all []coord, maxDistance int) int {
	distance := 0
	for _, a := range all {
		distance += manhattandistance(c, a)
		if distance > maxDistance {
			return distance
		}
	}
	return distance
}

func main2(all []coord, lx, hx, ly, hy int) {
	maxDistance := 10000
	count := 0
	for y := ly; y <= hy; y++ {
		for x := lx; x <= hx; x++ {
			if calcdistance(coord{x, y}, all, maxDistance) < maxDistance {
				count++
			}
		}
	}
	fmt.Println(count)
}

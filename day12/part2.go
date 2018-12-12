package main

import "fmt"

func main2(data []byte, offset int, rules map[byte]bool) {
	scanstart := 0

	maxiter := 50000000000

	deltas := make([]int64, 10)
	delta := int64(0)
	last := sumpots(data, offset)

	fmt.Printf("%05d %05d %010d ", 0, delta, last)
	prettyprintpots(data)

	for x := 0; x < maxiter; x++ {
		scanstart, data = itergroupins(data, rules, scanstart)

		t := sumpots(data, offset)
		fmt.Printf("%05d %05d %010d ", x+1, delta, t)
		prettyprintpots(data)

		delta = int64(t - last)
		last = t
		deltas[x%10] = delta
		if deltas[0] == deltas[1] &&
			deltas[1] == deltas[2] &&
			deltas[2] == deltas[3] &&
			deltas[3] == deltas[4] &&
			deltas[4] == deltas[5] &&
			deltas[5] == deltas[6] &&
			deltas[6] == deltas[7] &&
			deltas[7] == deltas[8] &&
			deltas[8] == deltas[9] {
			remaining := int64(maxiter-x) - 1
			fmt.Println("Stable state detected at iter ", x, " delta ", delta, " total ", last, " remaining ", remaining, " final score ", (last + (remaining * delta)))
			break
		}
	}

	fmt.Println(sumpots(data, offset))
}

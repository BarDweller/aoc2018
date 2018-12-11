package main

import (
	"fmt"
	"math"
)

//https://stackoverflow.com/a/46756052
func digit(num, place int) int {
	r := num % int(math.Pow(10, float64(place)))
	return r / int(math.Pow(10, float64(place-1)))
}

/*
Find the fuel cell's rack ID, which is its X coordinate plus 10.
Begin with a power level of the rack ID times the Y coordinate.
Increase the power level by the value of the grid serial number (your puzzle input).
Set the power level to itself multiplied by the rack ID.
Keep only the hundreds digit of the power level (so 12345 becomes 3; numbers with no hundreds digit become 0).
Subtract 5 from the power level.
*/
func cell(x, y, id int) (result int) {
	result = ((x + 10) * y) + id
	result *= (x + 10)
	result = digit(result, 3) - 5
	return result
}

func buildlayer1(id, max int) [][]int {
	result := make([][]int, max+1)
	for y := 1; y <= max; y++ {
		result[y] = make([]int, max+1)
		for x := 1; x <= max; x++ {
			result[y][x] = cell(x, y, id)
		}
	}
	return result
}

func buildlayer2(grid [][]int, max, side int) [][]int {
	result := make([][]int, max+1)
	for y := 1; y <= (max - (side - 1)); y++ {
		result[y] = make([]int, max+1)
		for x := 1; x <= (max - (side - 1)); x++ {
			for sx := 0; sx < side; sx++ {
				for sy := 0; sy < side; sy++ {
					result[y][x] += grid[y+sy][x+sx]
				}
			}
		}
	}
	return result
}

func seekmax(grid2 [][]int, max, side int) (tx, ty, tval int) {
	tval = 0
	for y := 1; y <= (max - (side - 1)); y++ {
		for x := 1; x <= (max - (side - 1)); x++ {
			if tval < grid2[y][x] {
				tval = grid2[y][x]
				tx = x
				ty = y
			}
		}
	}
	return tx, ty, tval
}

func main() {
	max := 300
	side := 3
	grid1 := buildlayer1(8141, max)
	grid := buildlayer2(grid1, max, side)

	fmt.Println(seekmax(grid, max, side))

	main2(grid1)
}

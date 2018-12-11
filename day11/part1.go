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

func calcsquare(grid [][]int, max, side int) (result int) {
	for sx := 1; sx <= side; sx++ {
		for sy := 1; sy <= side; sy++ {
			result += grid[sy][sx]
		}
	}
	return result
}

func buildlayer2(grid [][]int, max, side int) (maxx, maxy, power int) {
	result := make([][]int, max+1)
	willneverscorethis := max * max * 50
	firsty := willneverscorethis
	xsqr := 0
	for y := 1; y <= (max - (side - 1)); y++ {
		result[y] = make([]int, max+1)
		for x := 1; x <= (max - (side - 1)); x++ {
			if x == 1 {
				//first square in a row
				if firsty == willneverscorethis {
					//not just the first square in a row, also
					//the first square in a column =)
					//perform the first & only full square calc
					firsty = calcsquare(grid, max, side)
				} else {
					//move the square down a row
					for sx := 0; sx < side; sx++ {
						//remove old row
						firsty -= grid[y-1][x+sx]
						//add new row
						firsty += grid[y+side-1][x+sx]
					}
				}
				//clone square total to row total
				xsqr = firsty
				//save total into grid (for debug!)
				result[y][x] = firsty
			} else {
				//not the first square in a row.
				//move the square across a row
				for sy := 0; sy < side; sy++ {
					//remove old column
					xsqr -= grid[y+sy][x-1]
					//add new column
					xsqr += grid[y+sy][x+side-1]
				}
				//save total into grid (for debug!)
				result[y][x] = xsqr
			}
			//dow we have a new winner?
			if result[y][x] > power {
				maxx = x
				maxy = y
				power = result[y][x]
			}
		}
	}
	return maxx, maxy, power
}

func main() {
	max := 300
	side := 3
	grid := buildlayer1(8141, max)
	maxx, maxy, power := buildlayer2(grid, max, side)

	fmt.Println("X:", maxx, " Y:", maxy, " POWER:", power)
	main2(grid)
}

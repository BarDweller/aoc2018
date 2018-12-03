package main

import (
	"fmt"
	"strconv"
	"strings"
)

func loadData(input string) []string {
	lines := []string{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}

type claim struct {
	claimid          int
	xoffset, yoffset int
	width, height    int
}

func parseData(lines []string) []claim {
	claims := []claim{}
	for _, line := range lines {
		c := claim{}
		fmt.Sscanf(line, "#%d @ %d,%d: %dx%d", &c.claimid, &c.xoffset, &c.yoffset, &c.width, &c.height)
		claims = append(claims, c)
	}
	return claims
}

type coord struct {
	x, y int
}

type quilt struct {
	maxx, maxy int
	grid       map[coord]int
}

func assembleQuilt(claims []claim) quilt {
	grid := map[coord]int{}
	maxx, maxy := 0, 0
	for _, c := range claims {
		for y := c.yoffset; y < (c.yoffset + c.height); y++ {
			for x := c.xoffset; x < (c.xoffset + c.width); x++ {
				grid[coord{x: x, y: y}]++
				if x > maxx {
					maxx = x
				}
				if y > maxy {
					maxy = y
				}
			}
		}
	}
	return quilt{maxx: maxx, maxy: maxy, grid: grid}
}

func prettyPrintQuilt(q quilt) {
	for y := 0; y <= q.maxy; y++ {
		for x := 0; x <= q.maxx; x++ {
			val := q.grid[coord{x: x, y: y}]
			c := ""
			switch {
			case val == 0:
				{
					c = "."
				}
			case val > 0 && val < 10:
				{
					c = strconv.Itoa(val)
				}
			default:
				{
					c = "X"
				}
			}
			fmt.Print(c)
		}
		fmt.Println("")
	}
}

func countClaims(q quilt) int {
	count := 0
	for y := 0; y <= q.maxy; y++ {
		for x := 0; x <= q.maxx; x++ {
			if q.grid[coord{x: x, y: y}] > 1 {
				count++
			}
		}
	}
	return count
}

func main() {
	data := loadData(testdata())
	claims := parseData(data)
	quilt := assembleQuilt(claims)
	//prettyPrintQuilt(quilt)
	fmt.Println(countClaims(quilt))
	main2(quilt, claims)
}

func testdata() string {
	return `
	#1 @ 1,3: 4x4
	#2 @ 3,1: 4x4
	#3 @ 5,5: 2x2	
	`
}

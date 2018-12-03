package main

import "fmt"

func verifyClaims(q quilt, claims []claim) {
	for _, c := range claims {
		valid := true
		for y := c.yoffset; y < (c.yoffset + c.height); y++ {
			for x := c.xoffset; x < (c.xoffset + c.width); x++ {
				if q.grid[coord{x: x, y: y}] > 1 {
					valid = false
					break
				}
			}
			if !valid {
				break
			}
		}
		if valid {
			fmt.Println("Valid Claim: ", c)
		}
	}
}

func main2(q quilt, c []claim) {
	verifyClaims(q, c)
}

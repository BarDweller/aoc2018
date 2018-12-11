package main

import (
	"fmt"
)

type sideresult struct {
	side, x, y, maxpower int
}

func calcside(grid1 [][]int, max, side int, reply chan sideresult) {
	x, y, power := buildlayer2(grid1, max, side)
	reply <- sideresult{side: side, x: x, y: y, maxpower: power}
}

func main2(grid1 [][]int) {
	max := 300
	maxresult := sideresult{}
	var replies = make(chan sideresult)
	for side := max; side > 0; side-- {
		go calcside(grid1, max, side, replies)
	}
	replycount := 0
	for replycount < max {
		sideresult := <-replies
		replycount++
		fmt.Println("Side ", sideresult.side, " complete ", replycount, "/", max, " done.")
		if sideresult.maxpower > maxresult.maxpower {
			maxresult = sideresult
		}
	}
	fmt.Println("identifier : ", maxresult.x, ",", maxresult.y, ",", maxresult.side, "  with power ", maxresult.maxpower)
}

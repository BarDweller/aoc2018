package main

import (
	"fmt"
)

type node struct {
	prev  *node
	value int
	next  *node
}

func insert(newnode *node, currentnode *node) {
	//insert point is after the node after currentnode
	insertafter := currentnode.next

	newnode.next = insertafter.next
	newnode.prev = insertafter
	insertafter.next.prev = newnode
	insertafter.next = newnode
}

func remove(removalnode *node) {
	removalnode.next.prev = removalnode.prev
	removalnode.prev.next = removalnode.next
}

func doplacement(maxmarble int, playercount int) {
	playerscores := make([]int, playercount)

	zeromarblering := node{value: 0}
	zeromarblering.next = &zeromarblering
	zeromarblering.prev = &zeromarblering

	currentnode := &zeromarblering
	currentplayer := -1

	for marblenumber := 1; marblenumber <= maxmarble; marblenumber++ {
		currentplayer = (currentplayer + 1) % playercount
		if marblenumber%23 == 0 {
			playerscores[currentplayer] += marblenumber
			//removal point is the node 7 to the left of current node.
			removalnode := currentnode.prev.prev.prev.prev.prev.prev.prev
			currentnode = removalnode.next
			remove(removalnode)
			playerscores[currentplayer] += removalnode.value
		} else {
			marble := node{value: marblenumber}
			insert(&marble, currentnode)
			currentnode = &marble
		}
	}

	max := 0
	for _, val := range playerscores {
		if val > max {
			max = val
		}
	}
	fmt.Println("high score : ", max)
}

func main() {
	lastmarble := 72103
	players := 459
	doplacement(lastmarble, players)
	main2(lastmarble, players)
}

package main

import (
	"container/heap"
	"fmt"
)

//// BEGIN PQ CODE

//https://godoc.org/container/heap#ex-package--PriorityQueue
type Item struct {
	distance int //Item payload
	count    int //Item payload
	index    int //index of item, required by pq
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	//we want smaller priorties to be the 'best' option, as we're seeking smallest overall distance
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

///// END PQ CODE

func intmax(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func part2(input string) {
	nanos := loaddata(input)
	queue := PriorityQueue{}
	heap.Init(&queue)

	//build sorted queue of enter/leave overlap
	zero := coord{0, 0, 0}
	for _, nano := range nanos {
		distance := manhattandistance(zero, nano.loc)
		maxd := distance + nano.radius
		mind := distance - nano.radius
		heap.Push(&queue, &Item{distance: mind, count: 1})  //enter overlap
		heap.Push(&queue, &Item{distance: maxd, count: -1}) //leave overlap
	}

	//iterate queue finding maximum overlap zone
	chosendistance := 0
	overlaps := 0
	max := 0
	for len(queue) > 0 {
		item := (heap.Pop(&queue)).(*Item)
		overlaps += item.count
		if overlaps > max {
			max = overlaps
			chosendistance = item.distance
		}
	}
	fmt.Println("1d intersect distance ", chosendistance)
}

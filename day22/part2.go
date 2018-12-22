package main

import (
	"container/heap"
	"fmt"
)

const (
	none  = 0
	torch = 1
	rope  = 2
)

type step struct {
	loc  xy
	tool int
}

///// BEGIN PQ CODE

//https://godoc.org/container/heap#ex-package--PriorityQueue
type Item struct {
	loc      xy  //Item payload
	tool     int //Item payload
	priority int //used to sort items in queue, represents time in day22 puzzle
	index    int //index of item, required by pq
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	//we want smaller priorties to be the 'best' option, as we're seeking smallest overall time.
	return pq[i].priority < pq[j].priority
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

//get ALL possible neighbours for a point
func neighbours(loc xy) []xy {
	switch {
	case loc.x > 0 && loc.y > 0:
		return []xy{xy{loc.x + 1, loc.y}, xy{loc.x, loc.y + 1}, xy{loc.x - 1, loc.y}, xy{loc.x, loc.y - 1}}
	case loc.y > 0:
		return []xy{xy{loc.x + 1, loc.y}, xy{loc.x, loc.y + 1}, xy{loc.x, loc.y - 1}}
	case loc.x > 0:
		return []xy{xy{loc.x + 1, loc.y}, xy{loc.x, loc.y + 1}, xy{loc.x - 1, loc.y}}
	default:
		return []xy{xy{loc.x + 1, loc.y}, xy{loc.x, loc.y + 1}}
	}
}

//get where
func allowedneighbours(loc xy, tool int, g grid) []Item {
	var n []Item

	//abuse our xy type to create map from tool/target to other tool.
	othertool := map[xy]int{
		xy{0, 1}: 2,
		xy{1, 0}: 2,
		xy{0, 2}: 1,
		xy{2, 0}: 1,
		xy{1, 2}: 0,
		xy{2, 1}: 0,
	}

	for _, c := range neighbours(loc) {
		//get destination type.
		t := geterosion(c, g) % 3

		//destination / allowed tool
		//rock   == 0   rope | torch (2|1)
		//wet    == 1   rope | none  (2|0)
		//narrow == 2   none | torch (0|1)

		//so if our current tool doesn't match the destination type, we can move
		//there, at different costs, depending on if we change tool.

		if tool != t {
			//stick with current tool
			n = append(n, Item{loc: c, tool: tool, priority: 1})
			//swap to other tool
			n = append(n, Item{loc: c, tool: othertool[xy{tool, t}], priority: 8})
		}
	}

	return n
}

func part2(g grid) {

	queue := PriorityQueue{
		&Item{loc: xy{0, 0}, priority: 0, tool: torch},
	}
	heap.Init(&queue)

	//map of step to smallest time (priority) known for step.
	distances := map[step]int{}
	distances[step{xy{0, 0}, torch}] = 0

	for len(queue) > 0 {
		item := (heap.Pop(&queue)).(*Item)

		//if this item on the queue is the destination, AND we have the torch equipped,
		//then we're done here.
		if item.loc == g.target && item.tool == torch {
			fmt.Println("Time : ", item.priority)
			break
		}

		//have to allow the route to go past the target, but also have to
		//prevent it going off into infinity =)
		if item.loc.x > 5*g.target.x || item.loc.y > 5*g.target.y {
			//we're off in the weeds, give up!
			continue
		}

		//if we already know the current node, and we have a quicker route to it,
		//then skip processing this node.
		if time, found := distances[step{loc: item.loc, tool: item.tool}]; found && time < item.priority {
			continue
		}

		//for each possible (costed move to a) neighbour from the current location, add the move to the
		//queue if it represents a new cheaper option than we know already.
		for _, neighbour := range allowedneighbours(item.loc, item.tool, g) {
			nextstep := step{loc: neighbour.loc, tool: neighbour.tool}
			//add nodes we haven't seen before, or that have lower times than we know of for that
			//node so far.
			if time, found := distances[nextstep]; !found || item.priority+neighbour.priority < time {
				distances[nextstep] = item.priority + neighbour.priority
				heap.Push(&queue, &Item{loc: neighbour.loc, priority: item.priority + neighbour.priority, tool: neighbour.tool})
			}
		}
	}

}

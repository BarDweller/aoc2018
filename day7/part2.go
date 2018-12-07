package main

import (
	"fmt"
	"sort"
)

const processing int = 4

func walkNodes2(chosen node, nodemap *map[string]int, nodes map[string]*node) {

	//array of workers, with how long until they are free..

	//for example scenario
	timeTillWorkerFree := []int{0, 0}
	//for real data
	//timeTillWorkerFree := []int{0, 0, 0, 0, 0}

	//the node anygiven worker is working on.
	nodesByWorker := []string{"", "", "", "", ""}

	//total elapsed time so far.
	elapsedTime := 0

	availablenodes := []string{}
	for done := false; !done; done = len(availablenodes) == 0 {
		//mark the current node as 'processing', it will remain as such until enough
		//time elapses.
		(*nodemap)[chosen.label] = processing

		//figure out which nodes are 'available' for consideration from this node.
		for _, child := range chosen.children {
			if (*nodemap)[child.label] == unseen {
				(*nodemap)[child.label] = available
			}
		}
		availablenodes = []string{}
		for label, seen := range *nodemap {
			if seen == available {
				availablenodes = append(availablenodes, label)
			}
		}

		//no available nodes means we're all done =) the for loop above will exit us after this if stmt
		if len(availablenodes) > 0 {
			//sort candidates by alpha.
			sort.Strings(availablenodes)

			//first, check if there are any nodes possible..
			//(possible means candidate AND all parents are 'processed')
			nextNode := "UNKNOWN"
			workerIdx := -1
			for nextNode == "UNKNOWN" {

				for _, available := range availablenodes {
					allParentsFinished := true
					for _, parent := range nodes[available].parents {
						if parent.label != "" && (*nodemap)[parent.label] != processed {
							allParentsFinished = false
							break
						}
					}
					if allParentsFinished {
						nextNode = available
						break
					}
				}

				//if not.. then we need to advance time, and complete nodes, and retry
				//until we find that a node is now possible.

				//advance time by smallest amount required to unlock a worker (complete an item)
				minTime := 4096 //cheat, we'll never have a value this large
				for worker, time := range timeTillWorkerFree {
					//ignore workers that are idle, if we have no next node to give them yet.
					//instead we want to select the shortest time available to unlock the next
					//worker so we can recalculate the possible nodes.
					if time < minTime && !(time == 0 && nextNode == "UNKNOWN") {
						minTime = time
						workerIdx = worker
					}
				}

				//if minTime is >0 then all workers are busy, or we have no
				//possible next nodes because we need to complete ones in progress
				//we need to skip forward in time until the next worker frees up.
				//this will complete a node, and potentially unlock new candidate nodes.
				if minTime != 0 {
					elapsedTime += minTime
					for worker := range timeTillWorkerFree {
						if timeTillWorkerFree[worker] != 0 {
							timeTillWorkerFree[worker] -= minTime
						}
					}
				}
				//time just advanced by minTime, look at the workers who
				//were busy on something, and see if they just finished,
				//if so, move their node from processing to processed
				for idx, nodename := range nodesByWorker {
					if nodename != "" {
						if timeTillWorkerFree[idx] == 0 {
							(*nodemap)[nodesByWorker[idx]] = processed
						}
					}
				}
			}

			//so now we know which node we'll do next, and with which worker.
			if (*nodes[nextNode]).label != "" {
				//for example scenario...
				time := int(rune((*nodes[nextNode]).label[0]) - 'A' + 1)

				//for real data
				//time := int(rune((*nodes[nextNode]).label[0]) - 'A' + 61)

				//we know the workerIdx is at zero, so set it to the time required
				//to process this item
				timeTillWorkerFree[workerIdx] = time
				//assign this node to this worker, so we can mark it processed later.
				nodesByWorker[workerIdx] = (*nodes[nextNode]).label
			}

			//set our chosen node to be the one selected.
			chosen = *nodes[nextNode]
		}
	}

	//we're all out of 'new' nodes to process, but workers are still completing tasks.
	//find the longest worker, and add it to the current time.
	max := 0
	for _, time := range timeTillWorkerFree {
		if time > max {
			max = time
		}
	}
	elapsedTime += max

	fmt.Println("Elapsed time ", elapsedTime)
}

func main2(root node, nodes map[string]*node) {
	visitmap := &(map[string]int{})
	walkNodes2(root, visitmap, nodes)
}

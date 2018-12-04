package main

import (
	"fmt"
	"time"
)

func buildSleepMap(events []event) map[int]map[int]int {
	sleepPerMinutePerGuard := map[int]map[int]int{}
	sleepStart := time.Now()
	for _, event := range events {
		switch event.action {
		case sleep:
			sleepStart = event.date
		case wake:
			for minute := sleepStart.Minute(); minute < event.date.Minute(); minute++ {
				if sleepPerMinutePerGuard[minute] == nil {
					sleepPerMinutePerGuard[minute] = map[int]int{}
				}
				sleepPerMinutePerGuard[minute][event.guard]++
			}
		default:
		}
	}
	return sleepPerMinutePerGuard
}

func locateHighestCountedMinuteAndGuard(sleepMap map[int]map[int]int) (int, int) {

	max := 0
	chosenMinute := 0
	chosenGuard := 0

	for minute, guardMap := range sleepMap {
		for guard, count := range guardMap {
			if count > max {
				max = count
				chosenMinute = minute
				chosenGuard = guard
			}
		}
	}

	return chosenMinute, chosenGuard
}

func main2(events []event) {
	sleepMap := buildSleepMap(events)
	chosenMinute, chosenGuard := locateHighestCountedMinuteAndGuard(sleepMap)
	fmt.Printf("Part2: Chose Minute %d and Guard %d with result %d\n", chosenMinute, chosenGuard, (chosenMinute * chosenGuard))
}

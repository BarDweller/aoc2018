package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
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

type action int

const (
	begin action = 0
	sleep action = 1
	wake  action = 2
)

type event struct {
	date   time.Time
	guard  int
	action action
}

func (a action) String() string {
	switch a {
	case begin:
		return "begin"
	case sleep:
		return "sleep"
	case wake:
		return "wake"
	}
	panic("Unknown enum value")
}

func (e event) String() string {

	return e.date.Format("[2006-01-02 15:04] ") + strconv.Itoa(e.guard) + " " + e.action.String() + "\n"
}

func parseData(lines []string) []event {
	layout := "[2006-01-02 15:04]"
	events := []event{}
	for _, line := range lines {
		e := event{}
		parts := strings.SplitAfter(line, "]")
		date, err := time.Parse(layout, parts[0])
		if err != nil {
			panic("Could not parse date in line " + parts[0])
		}
		e.date = date
		switch parts[1] {
		case " falls asleep":
			{
				e.action = sleep
			}
		case " wakes up":
			{
				e.action = wake
			}
		default:
			{
				e.action = begin
				fmt.Sscanf(parts[1], " Guard #%d begins shift", &e.guard)
			}
		}
		events = append(events, e)
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].date.Before(events[j].date)
	})
	guard := -1
	for idx := range events {
		switch events[idx].action {
		case begin:
			{
				guard = events[idx].guard
			}
		default:
			{
				events[idx].guard = guard
			}
		}
	}
	return events
}

func calcMinutesAsleep(events []event) map[int]int {
	sleepPerGuard := map[int]int{}
	sleepStart := time.Now()
	for _, event := range events {
		switch event.action {
		case sleep:
			sleepStart = event.date
		case wake:
			duration := event.date.Sub(sleepStart)
			sleepPerGuard[event.guard] += int(duration.Minutes())
		default:
		}
	}
	return sleepPerGuard
}

func selectGuardWithMostMinutesAsleep(sleepMap map[int]int) int {
	max := 0
	chosenGuard := 0
	for k, v := range sleepMap {
		if v > max {
			chosenGuard = k
			max = v
		}
	}
	return chosenGuard
}

func findMostOftenAlseepMinute(events []event, guard int) int {
	sleepPerMinute := map[int]int{}
	sleepStart := time.Now()
	for _, event := range events {
		if event.guard == guard {
			switch event.action {
			case sleep:
				sleepStart = event.date
			case wake:
				for minute := sleepStart.Minute(); minute < event.date.Minute(); minute++ {
					sleepPerMinute[minute]++
				}
			default:
			}
		}
	}
	max := 0
	chosenMinute := 0
	for k, v := range sleepPerMinute {
		if v > max {
			chosenMinute = k
			max = v
		}
	}
	return chosenMinute
}

func main() {
	data := loadData(testdata())
	events := parseData(data)
	sleepPerGuard := calcMinutesAsleep(events)
	chosenGuard := selectGuardWithMostMinutesAsleep(sleepPerGuard)
	chosenMinute := findMostOftenAlseepMinute(events, chosenGuard)
	fmt.Printf("Part1: Chose Minute %d and Guard %d with result %d\n",chosenMinute, chosenGuard, (chosenGuard * chosenMinute))
	main2(events)
}

func testdata() string {
	return `
	[1518-11-05 00:03] Guard #99 begins shift	
	[1518-11-01 00:00] Guard #10 begins shift
	[1518-11-01 00:05] falls asleep
	[1518-11-01 00:25] wakes up
	[1518-11-01 00:30] falls asleep
	[1518-11-01 00:55] wakes up
	[1518-11-01 23:58] Guard #99 begins shift
	[1518-11-02 00:40] falls asleep
	[1518-11-02 00:50] wakes up
	[1518-11-03 00:05] Guard #10 begins shift
	[1518-11-03 00:24] falls asleep
	[1518-11-03 00:29] wakes up
	[1518-11-04 00:02] Guard #99 begins shift
	[1518-11-04 00:36] falls asleep
	[1518-11-04 00:46] wakes up
	[1518-11-05 00:45] falls asleep
	[1518-11-05 00:55] wakes up
	`
}

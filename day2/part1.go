package main

import (
	"fmt"
	"sort"
	"strings"
)

func processLine(line string) (int, int) {
	twos := 0
	threes := 0

	//split string into slice of strings, and sort it
	letters := strings.Split(line, "")
	sort.Strings(letters)

	//add a terminator ;)
	letters = append(letters, "*")

	//look for contiguous sequences, no end case needed
	//due to terminator char =)
	last := letters[0]
	count := 1
	for _, letter := range letters[1:len(letters)] {
		if letter == last {
			count++
		} else {
			//sequence is broken, how long was it?
			if count == 2 {
				twos++
			}
			if count == 3 {
				threes++
			}
			//reset for next
			last = letter
			count = 1
		}
	}
	return twos, threes
}

func main() {
	twos := 0
	threes := 0
	for _, line := range strings.Split(testdata(), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		linetwos, linethrees := processLine(line)

		if linetwos > 0 {
			twos++
		}
		if linethrees > 0 {
			threes++
		}
	}
	fmt.Println("Part 1 Result ", twos, threes, (twos * threes))
	main2()
}

func testdata() string {
	return `
	abcdef
	bababc
	abbcde
	abcccd
	aabcdd
	abcdee
	ababab
	`
}

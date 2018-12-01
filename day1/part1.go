package main

import (
	"fmt"
	"strconv"
	"strings"
)

func asInt(txt string) int {
	if val, err := strconv.Atoi(txt); err == nil {
		return val
	} else {
		panic("bad input, expected number, got : " + txt)
	}
}

func processLine(line string, currentTotal int) int {
	line = strings.TrimSpace(line)
	if line != "" {
		if strings.HasPrefix(line, "+") {
			line = line[1:len(line)]
		}
		val := asInt(line)
		currentTotal += val
	}
	return currentTotal
}

func main() {
	total := 0
	for _, line := range strings.Split(testdata(), "\n") {
		total = processLine(line, total)
	}
	fmt.Println("Part 1 Result ", total)
	main2()
}

func testdata() string {
	return `
	+1
	+1
	-1
	`
}

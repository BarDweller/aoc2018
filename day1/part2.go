package main

import (
	"fmt"
	"strings"
)

func main2() {
	total := 0
	seen := map[int]bool{}
	exit := false

	for !exit {
		for _, line := range strings.Split(testdata(), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			total = processLine(line, total)
			if seen[total] {
				exit = true
				fmt.Println("Repeated Total was ", total)
				break
			} else {
				seen[total] = true
			}
		}
	}
}

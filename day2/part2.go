package main

import (
	"fmt"
	"strings"
)

func loadData(data string) []string {
	lines := []string{}
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}

func compare(a string, b string) bool {
	mismatch := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			mismatch++
		}
		if mismatch > 1 {
			break
		}
	}
	return mismatch == 1
}

func combine(a string, b string) string {
	result := ""
	for i := 0; i < len(a); i++ {
		if a[i] == b[i] {
			result += string(a[i])
		}
	}
	return result
}

func main2() {
	lines := loadData(part2testdata())
	for idx, line := range lines {
		for idx2, otherline := range lines[idx:len(lines)] {
			if idx2 == 0 {
				continue
			}
			if compare(line, otherline) {
				fmt.Println("Hit: ", combine(line, otherline))
			}
		}
	}
}

func part2testdata() string {
	return `
	abcde
	fghij
	klmno
	pqrst
	fguij
	axcye
	wvxyz	
	`
}

package main

import (
	"fmt"
	"strings"
)

func dropChars(lc rune, uc rune, input string) string {
	drop := func(r rune) rune {
		if r == lc || r == uc {
			return -1
		} else {
			return r
		}
	}
	return strings.Map(drop, input)
}

func main2() {
	min := len(testdata())
	for lc, uc := 'a', 'A'; lc <= 'z'; lc, uc = lc+1, uc+1 {
		polymer := []rune(dropChars(lc, uc, testdata()))
		result := reactUntilComplete(polymer)
		if result < min {
			min = result
		}
	}
	fmt.Println("Shortest ", min)
}

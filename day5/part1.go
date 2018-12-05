package main

import (
	"fmt"
)

func react(polymer []rune, offset int) (bool, int, []rune) {
	last := polymer[offset]
	for idx := (offset + 1); idx < len(polymer); idx++ {
		current := polymer[idx]
		difference := 0
		if current > last {
			difference = int(current) - int(last)
		} else {
			difference = int(last) - int(current)
		}
		if difference == 32 {
			result := append(polymer[0:idx-1], polymer[idx+1:len(polymer)]...)
			if idx > 2 {
				return false, idx - 2, result
			} else {
				return false, 0, result
			}
		}
		last = current
	}
	return true, 0, polymer
}

func reactUntilComplete(polymer []rune) int {
	offset := 0
	done := false
	for !done {
		done, offset, polymer = react(polymer, offset)
	}
	return len(polymer)
}

func main() {
	fmt.Println("Count ", reactUntilComplete([]rune(testdata())))
	main2()
}

func testdata() string {
	return `dabAcCaCBAcCcaDA`
}

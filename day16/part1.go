package main

import (
	"fmt"
	"strings"
)

type regs [4]int

type state struct {
	beforeregs regs
	afterregs  regs
	opcode     int
	inputreg1  int
	inputreg2  int
	outputreg  int
}

type operation struct {
	name string
	fn   func(beforestate regs, A int, B int, C int) regs
}

func newstate() state {
	var before regs
	var after regs
	return state{before, after, 0, 0, 0, 0}
}

func loaddata(input string) []state {
	result := []state{}

	s := newstate()
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		switch {
		case strings.Contains(line, "Before:"):
			fmt.Sscanf(line, "Before: [%d, %d, %d, %d]",
				&s.beforeregs[0],
				&s.beforeregs[1],
				&s.beforeregs[2],
				&s.beforeregs[3])
		case strings.Contains(line, "After:"):
			fmt.Sscanf(line, "After: [%d, %d, %d, %d]",
				&s.afterregs[0],
				&s.afterregs[1],
				&s.afterregs[2],
				&s.afterregs[3])
			//store completed s
			result = append(result, s)
			//create next s
			s = newstate()
		default:
			fmt.Sscanf(line, "%d %d %d %d]",
				&s.opcode,
				&s.inputreg1,
				&s.inputreg2,
				&s.outputreg)
		}
	}
	return result
}

func addr(r regs, A, B, C int) regs {
	var result = r
	result[C] = r[A] + r[B]
	return result
}

func addi(r regs, A, B, C int) regs {
	var result = r
	result[C] = r[A] + B
	return result
}

func mulr(r regs, A, B, C int) regs {
	var result = r
	result[C] = r[A] * r[B]
	return result
}

func muli(r regs, A, B, C int) regs {
	var result = r
	result[C] = r[A] * B
	return result
}

func banr(r regs, A, B, C int) regs {
	var result = r
	result[C] = r[A] & r[B]
	return result
}

func bani(r regs, A, B, C int) regs {
	var result = r
	result[C] = r[A] & B
	return result
}

func borr(r regs, A, B, C int) regs {
	var result = r
	result[C] = r[A] | r[B]
	return result
}

func bori(r regs, A, B, C int) regs {
	var result = r
	result[C] = r[A] | B
	return result
}

func setr(r regs, A, B, C int) regs {
	var result = r
	result[C] = r[A]
	return result
}

func seti(r regs, A, B, C int) regs {
	var result = r
	result[C] = A
	return result
}

func gtir(r regs, A, B, C int) regs {
	var result = r
	result[C] = 0
	if A > r[B] {
		result[C] = 1
	}
	return result
}

func gtri(r regs, A, B, C int) regs {
	var result = r
	result[C] = 0
	if r[A] > B {
		result[C] = 1
	}
	return result
}

func gtrr(r regs, A, B, C int) regs {
	var result = r
	result[C] = 0
	if r[A] > r[B] {
		result[C] = 1
	}
	return result
}
func eqir(r regs, A, B, C int) regs {
	var result = r
	result[C] = 0
	if A == r[B] {
		result[C] = 1
	}
	return result
}

func eqri(r regs, A, B, C int) regs {
	var result = r
	result[C] = 0
	if r[A] == B {
		result[C] = 1
	}
	return result
}

func eqrr(r regs, A, B, C int) regs {
	var result = r
	result[C] = 0
	if r[A] == r[B] {
		result[C] = 1
	}
	return result
}

func commandtable() []operation {
	r := []operation{}
	r = append(r, operation{"addr", addr})
	r = append(r, operation{"addi", addi})
	r = append(r, operation{"mulr", mulr})
	r = append(r, operation{"muli", muli})
	r = append(r, operation{"banr", banr})
	r = append(r, operation{"bani", bani})
	r = append(r, operation{"borr", borr})
	r = append(r, operation{"bori", bori})
	r = append(r, operation{"setr", setr})
	r = append(r, operation{"seti", seti})
	r = append(r, operation{"gtir", gtir})
	r = append(r, operation{"gtri", gtri})
	r = append(r, operation{"gtrr", gtrr})
	r = append(r, operation{"eqir", eqir})
	r = append(r, operation{"eqri", eqri})
	r = append(r, operation{"eqrr", eqrr})
	return r
}

func countMatches(s state, ops []operation, results chan int) {
	count := 0
	for _, op := range ops {
		result := op.fn(s.beforeregs, s.inputreg1, s.inputreg2, s.outputreg)
		if result == s.afterregs {
			count++
		}
	}
	results <- count
}

func part1(samples []state, ops []operation) {
	results := make(chan int)
	for _, s := range samples {
		go countMatches(s, ops, results)
	}
	count := 0
	for range samples {
		sampleCount := <-results
		if sampleCount >= 3 {
			count++
		}
	}
	fmt.Println("Count", count)
}

func main() {
	s := loaddata(testdata())
	ops := commandtable()
	part1(s, ops)
	main2(s, ops)
}

func testdata() string {
	return `
	Before: [3, 2, 1, 1]
	9 2 1 2
	After:  [3, 2, 2, 1]
	Before: [0, 2, 2, 2]
	6 0 0 3
	After:  [0, 2, 2, 0]	
	`
}

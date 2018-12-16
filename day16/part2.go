package main

import (
	"fmt"
	"strings"
)

type opmap map[int]int

type opcoderesult struct {
	opcode  int
	matches map[int]bool
}

func getMatchIndexes(s state, ops []operation, om opmap, results chan opcoderesult) {
	count := map[int]bool{}
	for i, op := range ops {
		result := op.fn(s.beforeregs, s.inputreg1, s.inputreg2, s.outputreg)
		if result == s.afterregs {
			count[i] = true
		}
	}
	results <- opcoderesult{s.opcode, count}
}

func resolveops(samples []state, ops []operation) opmap {
	om := opmap{}
	results := make(chan opcoderesult)
	launched := 0

	//evaluate all the samples in parallel
	for _, s := range samples {
		go getMatchIndexes(s, ops, om, results)
		launched++
	}

	//aggregate results by opcode driven
	mergedresults := map[int]map[int]bool{}
	for ; launched > 0; launched-- {
		//accept result from channel
		matchresult := <-results
		//if we haven't seen this opcode before, we'll need to create the map
		if _, found := mergedresults[matchresult.opcode]; !found {
			mergedresults[matchresult.opcode] = map[int]bool{}
		}
		//for each key/value in the match result map, add to
		//results for the opcode that was tested.
		for k, v := range matchresult.matches {
			mergedresults[matchresult.opcode][k] = v
		}
	}

	//now identify which ops are which indexes!
	//keep going until we have a match for every opcode
	for len(om) < len(ops) {
		//remove identified opcodes for this pass
		for _, matchingindexes := range mergedresults {
			for _, index := range om {
				delete(matchingindexes, index)
			}
		}
		//look for any opcodes that now only have a single index
		for opcode, matchingindexes := range mergedresults {
			if len(matchingindexes) == 1 {
				for k := range matchingindexes {
					om[opcode] = k
				}
			}
		}
	}
	return om
}

func runprog(input string, om opmap, ops []operation) {
	var r regs
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		//read in next instruction
		opcode, A, B, C := 0, 0, 0, 0
		fmt.Sscanf(line, "%d %d %d %d]",
			&opcode,
			&A,
			&B,
			&C)

		//run instruction, using opmap to locate appropriate fn to drive
		r = ops[om[opcode]].fn(r, A, B, C)
	}

	//dump answer for part 2.
	fmt.Println("Value in Reg 0 at end of prog ", r[0])
}

func main2(s []state, ops []operation) {
	om := resolveops(s, ops)
	runprog(testprog(), om, ops)
}

func testprog() string {
  return ``
}

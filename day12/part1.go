package main

import (
	"fmt"
	"strings"
)

func appendbit(data []byte, idx int, on bool) {
	byteidx := (uint)(idx / 8)
	bitidx := (uint)(idx % 8)
	if on {
		data[byteidx] |= 1 << (7 - bitidx)
	} else {
		data[byteidx] &^= 1 << (7 - bitidx)
	}
}

func parseinput(input string, padbytes int, postpad int) []byte {

	//allow extra space for future expansions..
	emptybyte := "........"
	pad := ""
	ppad := ""

	for i := 0; i < padbytes; i++ {
		pad += emptybyte
	}

	for i := 0; i < postpad; i++ {
		ppad += emptybyte
	}

	input = pad + input + ppad

	data := make([]byte, (len(input)/8)+1)
	for i, c := range input {
		switch c {
		case '#':
			appendbit(data, i, true)
		case '.':
			appendbit(data, i, false)
		}
	}

	return data
}

func itergroupins(data []byte, rules map[byte]bool, from int) (int, []byte) {
	result := make([]byte, len(data))
	firstdata := 0
	datastarted := false
	last := 0
	for i := from; i < len(data); i++ {
		b := data[i]
		//skip all the leading pad zeros added for negative growth.
		if b != 0 && !datastarted {
			datastarted = true
			firstdata = i - 1
		}
		if datastarted && (b != 0 || data[i-1] != 0) {
			last = i
			p := uint16(data[i-1])
			c := uint16(data[i])
			t := (p << 8) + c
			//got 16 bits, move 5 bit window across, isolate patterns
			for x := 0; x < 8; x++ {
				m := t & 0xF800
				s := m >> 11
				q := byte(s)
				//todo, why -6 !!
				appendbit(result, (i*8)+x-6, rules[q])
				t = t << 1
			}
		}
	}
	if last+1 == len(data) {
		result = append(result, 0)
	}
	return firstdata, result
}

func evalpattern(p byte) bool {
	return true
}

func sumpots(data []byte, offset int) int64 {
	total := int64(0)
	for i, b := range data {
		for bit := 0; bit < 8; bit++ {
			potnumber := int64(-offset*8) + int64(i*8) + int64(bit)
			if (b<<uint(bit))&0x80 == 0x80 {
				total += int64(potnumber)
			}
		}
	}
	return total
}

func prettyprintpots(data []byte) {
	for _, b := range data {
		for bit := 0; bit < 8; bit++ {
			if (b<<uint(bit))&0x80 == 0x80 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}

		}
	}
	fmt.Println("")
}

func processrules(input string) map[byte]bool {
	result := map[byte]bool{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = "..." + line
		if line[12:13] == "#" {
			data := []byte{0}
			for i, r := range line[0:8] {
				switch r {
				case '#':
					appendbit(data, i, true)
				case '.':
					appendbit(data, i, false)
				}
			}
			result[data[0]] = true
		}
	}
	return result
}

func twentyiters(data []byte, offset int, rules map[byte]bool) {
	scanstart := 0
	for x := 0; x < 20; x++ {
		scanstart, data = itergroupins(data, rules, scanstart)
	}
	fmt.Println(sumpots(data, offset))
}

func main() {

	//real input.
	//input := "..##.#######...##.###...#..#.#.#..#.##.#.##....####..........#..#.######..####.#.#..###.##..##..#..#"
	//rules := processrules(rules())

	//test input
	input := "#..#.#..##......###...###"
	rules := processrules(testrules())

	//offset to zero bit byte / padding bytes each end.
	offset := 2
	data := parseinput(input, offset, offset)

	twentyiters(data, offset, rules)

	main2(data, offset, rules)
}

func testrules() string {
	return `
	...## => #
	..#.. => #
	.#... => #
	.#.#. => #
	.#.## => #
	.##.. => #
	.#### => #
	#.#.# => #
	#.### => #
	##.#. => #
	##.## => #
	###.. => #
	###.# => #
	####. => #	
	`
}

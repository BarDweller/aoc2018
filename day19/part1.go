package main

import (
	"fmt"
	"strings"
)

type regs [6]int

type operation struct {
	name string
	fn   func(beforestate regs, A int, B int, C int) regs
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

type commandmap map[string]operation

func buildcommandmap(ops []operation) commandmap {
	cm := map[string]operation{}
	for _, c := range ops {
		cm[c.name] = c
	}
	return cm
}

type instruction struct {
	op      string
	A, B, C int
}

type program struct {
	ip         int
	statements []instruction
}

func loaddata(input string) program {
	p := program{0, []instruction{}}
	ipset := false
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !ipset && line[0] != '#' {
			panic("Bad input, expected #ip at line 0")
		}
		if ipset {
			i := instruction{}
			fmt.Sscanf(line, "%s %d %d %d", &i.op, &i.A, &i.B, &i.C)
			p.statements = append(p.statements, i)
		} else {
			fmt.Sscanf(line, "#ip %d", &p.ip)
			ipset = true
		}
	}
	return p
}

type machinestate struct {
	r regs
}

func runprogram(p program, ms machinestate, cm commandmap) machinestate {
	ipidx := p.ip
	ip := 0
	for ip >= 0 && ip < len(p.statements) {
		s := p.statements[ip]
		ms.r[ipidx] = ip
		//fmt.Print("ip=", ip, ms.r, " ", s.op, " ", s.A, s.B, s.C, " ")
		ms.r = cm[s.op].fn(ms.r, s.A, s.B, s.C)
		//fmt.Println(ms.r)
		ip = ms.r[ipidx]
		ip++
	}
	return ms
}

func main() {

	cm := buildcommandmap(commandtable())
	ms := machinestate{}

	p := loaddata(testdata())

	ms = runprogram(p, ms, cm)

	fmt.Println("Part1 Reg 0 ", ms.r[0])

	part2()
}

func testdata() string {
	return `
	#ip 0
	seti 5 0 1
	seti 6 0 2
	addi 0 1 0
	addr 1 2 3
	setr 1 0 0
	seti 8 0 4
	seti 9 0 5	
	`
}

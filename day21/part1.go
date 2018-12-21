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

		//the only compare inst that can break the loop is at r28, which is ip=29
		//at this point in the code, exit as soon as we hit this point, the value
		//being compared against is the lowest answer we can exit with.
		if ip == 29 {
			break
		}
	}
	return ms
}

func main() {

	cm := buildcommandmap(commandtable())
	ms := machinestate{}

	p := loaddata(data())

	ms = runprogram(p, ms, cm)

	//dump r4, because eqrr 4 0 2 is the only way it can break the loop.
	//so if we know r4, we know the first value it can exit at, which 
	//is the "lowest non-negative integer value for register 0 that causes the program to halt after executing the fewest instructions"
	fmt.Println("Part1 Reg 4 ", ms.r[4])

	part2()
}
func data() string {
	return `
	#ip 5
	seti 123 0 4
	bani 4 456 4
	eqri 4 72 4
	addr 4 5 5
	seti 0 0 5
	seti 0 8 4
	bori 4 65536 3
	seti 14464005 5 4
	bani 3 255 2
	addr 4 2 4
	bani 4 16777215 4
	muli 4 65899 4
	bani 4 16777215 4
	gtir 256 3 2
	addr 2 5 5
	addi 5 1 5
	seti 27 7 5
	seti 0 3 2
	addi 2 1 1
	muli 1 256 1
	gtrr 1 3 1
	addr 1 5 5
	addi 5 1 5
	seti 25 2 5
	addi 2 1 2
	seti 17 9 5
	setr 2 2 3
	seti 7 3 5
	eqrr 4 0 2
	addr 2 5 5
	seti 5 9 5	
	`
}

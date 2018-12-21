package main

import "fmt"

func runprogram2(p program, ms machinestate, cm commandmap) int {
	ipidx := p.ip
	ip := 0
	seen := map[int]bool{}
	last := 0
	for ip >= 0 && ip < len(p.statements) {
		s := p.statements[ip]
		ms.r[ipidx] = ip
		//fmt.Print("ip=", ip, ms.r, " ", s.op, " ", s.A, s.B, s.C, " ")
		ms.r = cm[s.op].fn(ms.r, s.A, s.B, s.C)
		//fmt.Println(ms.r)
		ip = ms.r[ipidx]
		ip++

		if ip == 29 {
			//we're at the compare instruction for the exit,
			//track the set of values that can cause exit, 
			//the last one seen before the set begins to repeat
			//will be "the lowest non-negative integer value for register 0 that causes the program to halt after executing the most instructions"
			if _, found := seen[ms.r[4]]; found {
				break
			}
			seen[ms.r[4]] = true
			last = ms.r[4]
		}
	}
	return last
}

func part2() {
	cm := buildcommandmap(commandtable())
	ms := machinestate{}

	p := loaddata(data())

	last := runprogram2(p, ms, cm)

	fmt.Println("Part2 Last val before repeat was ", last)
}

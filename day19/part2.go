package main

import "fmt"

func part2() {
	r0, r4, r5 := 0, 0, 0

	// addi 2 16 2  ip+=16 //goto l17

	//l1: //OUTER LOOP:
	// seti 1 0 1  r1 = 1
	//l2: //INNER LOOP:
	// seti 1 3 3  r3 = 1
	// mulr 1 3 5  r5 = r1 * r3
	// eqrr 5 4 5  test r5==r4, result in r5
	// addr 5 2 2  skip next inst if r5 true
	// addi 2 1 2  skip next inst.
	// addr 1 0 0  if r5==r4 add 1 to r0
	// addi 3 1 3  r3 += 1
	// gtrr 3 4 5  test r3==r4, result in r5
	// addr 2 5 2  skip next inst if r5 true
	// seti 2 6 2  ip = 2 (INNER LOOP)
	// addi 1 1 1  r1+=1
	// gtrr 1 4 5  test r1==r4, result in r5
	// addr 5 2 2  skip next inst if r5 true.
	// seti 1 1 2 ip = 1 (OUTER LOOP)

	// mulr 2 2 2 ip = ip * ip // This triggers the final code exit.

	//basically
	//for r1 := 1 to r4 {
	//  for r3 := 1 to r4 {
	//    r5 = r1 * r3
	//    if r5 == r4 { r0+=r1 }
	//  }
	//}
	//
	//eg.
	// loop thru 0..n * 0..n and add n if the product is r4
	//or.
	// sum of everything that fits cleanly into r4, n times.
	//or.
	// sum all the divisors of r4

	//part1 (from setup, below)
	r4 = 2 * 2 * 19 * 11
	r5 = (6 * 22) + 19
	r4 += r5

	//part2 (from setup, below)
	r5 = ((27 * 28) + 29) * 30 * 14 * 32
	r4 += r5

	//sum the divisors of r4
	for t := 1; t <= r4; t++ {
		if r4%t == 0 {
			r0 += t
		}
	}

	fmt.Println("Part2 Reg 0 ", r0)

	//setup...
	//l17:
	// addi 4 2 4  r4 = 2 (assuming 0 start)
	// mulr 4 4 4  r4 = 4
	// mulr 2 4 4  r4 = 4 * 19 (ip=19)
	// muli 4 11 4 r4 = 4 * 19 * 11 == 836

	// addi 5 6 5  r5 = 6 (assuming 0 start)
	// mulr 5 2 5  r5 = 6 * 22 (ip==22)
	// addi 5 19 5 r5 = 6 * 22 + 19
	// addr 4 5 4  r4 = (6 * 22 + 19) + 836 == 987

	// addr 2 0 2  skip next bit if r0 is 0
	// seti 0 7 2  goto l1  //EXIT setup HERE if Part 1

	// setr 2 6 5  r5 = 27 (ip==27)
	// mulr 5 2 5  r5 = 27 * 28 (ip==28)
	// addr 2 5 5  r5 = (27 * 28) + 29 (ip==29)
	// mulr 2 5 5  r5 = ((27 * 28) + 29) * 30 == 23550 (ip==30)
	// muli 5 14 5 r5 = 23550 * 14
	// mulr 5 2 5  r5 = 23550 * 14 * 32 == 10550400
	// addr 4 5 4  r4 += r5 (r4==987, r5==10550400)
  
	// seti 0 7 0  r0=0 // no-op
	// seti 0 3 2  goto l1
}

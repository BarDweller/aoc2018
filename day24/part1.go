package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type group struct {
	side        int
	units       int
	hp          int
	weaknesses  map[string]bool
	immunities  map[string]bool
	attackpower int
	attacktype  string
	initiative  int
}

func loaddata(input string) []group {
	gs := []group{}
	re := regexp.MustCompile(`(\d+) units each with (\d+) hit points \(([^)]+)\) with an attack that does (\d+) ([^ ]+) damage at initiative (\d+)`)
	state := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if state == 0 && line == "Immune System:" {
			state = 1
			continue
		}
		if state == 1 && line == "Infection:" {
			state = 2
			continue
		}
		if state == 1 || state == 2 {
			g := group{weaknesses: map[string]bool{}, immunities: map[string]bool{}}
			res := re.FindStringSubmatch(line)
			for i, v := range res {
				switch i {
				case 1:
					if u, err := strconv.Atoi(v); err == nil {
						g.units = u
					}
				case 2:
					if u, err := strconv.Atoi(v); err == nil {
						g.hp = u
					}
				case 3:
					v = strings.TrimSpace(v)
					if v == "" {
						continue
					}
					for _, p := range strings.Split(v, ";") {
						if strings.Contains(p, "weak to") {
							for _, w := range strings.Split(p[len("weak to "):], ", ") {
								w = strings.TrimSpace(w)
								g.weaknesses[w] = true
							}
						} else {
							for _, i := range strings.Split(p[len("immune to "):], ", ") {
								i = strings.TrimSpace(i)
								g.immunities[i] = true
							}
						}
					}
				case 4:
					if u, err := strconv.Atoi(v); err == nil {
						g.attackpower = u
					}
				case 5:
					g.attacktype = strings.TrimSpace(v)
				case 6:
					if u, err := strconv.Atoi(v); err == nil {
						g.initiative = u
					}
				}
			}
			g.side = state
			if g.initiative == 0 {
				fmt.Println(line)
			}
			gs = append(gs, g)
		}
	}
	return gs
}

func warisover(gs []group) bool {
	sides := map[int]int{}
	for _, g := range gs {
		if g.units > 0 {
			sides[g.side]++
		}
	}
	return sides[1] == 0 || sides[2] == 0
}

func calcdmg(attacker group, defender group) int {
	basedmg := attacker.units * attacker.attackpower
	dmg := basedmg
	if defender.weaknesses[attacker.attacktype] {
		dmg *= 2
	}
	if defender.immunities[attacker.attacktype] {
		dmg = 0
	}
	return dmg
}

func part1(r []group) int {

	//sanity test.. can we use initiative as a unique key ?
	test := map[int]int{}
	for i, g := range r {
		if _, found := test[g.initiative]; found {
			fmt.Println(g)
			panic("Non unique initiative found")
		} else {
			test[g.initiative] = i
		}
	}
	//are all the keys sequential ?
	for i := 1; i < len(test); i++ {
		if _, found := test[i]; !found {
			panic("initiatives are not sequential")
		}
	}

	//create index by initiative
	byinitiative := []*group{}
	//add zero index empty group.
	byinitiative = append(byinitiative, &group{})
	//append rest in order.
	for i := 1; i <= len(test); i++ {
		byinitiative = append(byinitiative, &r[test[i]])
	}

	for !warisover(r) {

		//sort by effective power
		sort.Slice(r, func(i, j int) bool {
			iep := (r[i].units * r[i].attackpower)
			jep := (r[j].units * r[j].attackpower)
			if iep == jep {
				return r[i].initiative > r[j].initiative
			} else {
				return iep > jep
			}
		})

		//gah.. sort rewrites the underlying array, which ruins
		//the byinitiative map. #lessonlearned
		//rebuild it.
		for i := 0; i < len(r); i++ {
			byinitiative[r[i].initiative] = &r[i]
		}

		//allocate targets

		//initiative(defender) -> attacker selected
		targetallocated := map[int]bool{}
		//intiative(attacker) -> defender selected
		targetselection := map[int]*group{}

		for _, attacker := range r {
			//dead units attack no one!
			if attacker.units <= 0 {
				continue
			}
			max := 0
			selected := -1
			for _, defender := range r {
				//only attack other side (also skips self ;) )
				if attacker.side == defender.side {
					continue
				}
				//don't attack dead people, it's not nice.
				if defender.units <= 0 {
					continue
				}
				//don't attack someone already selected
				if _, found := targetallocated[defender.initiative]; found {
					continue
				}

				//evaluate dmg to deal
				dmg := calcdmg(attacker, defender)

				//fmt.Println("attacker (", attacker.initiative, ") would deal ", dmg, " to defending group (", defender.initiative, ")")
				if dmg > max {
					selected = defender.initiative
					max = dmg
				}
			}
			if selected != -1 {
				//fmt.Println("attacker (", attacker.initiative, ") has selected (", selected, ") as its target")
				targetallocated[selected] = true
				targetselection[attacker.initiative] = byinitiative[selected]
			}
		}

		//fmt.Println(r)

		//attack phase
		killed := 0
		for i := len(byinitiative) - 1; i > 0; i-- {
			if target, found := targetselection[i]; found {
				attacker := byinitiative[i]
				if attacker.units <= 0 {
					continue
				}
				//have to recalculate dmg, because unit counts will have been affected by prior attacks.
				dmg := calcdmg(*attacker, *target)
				//fmt.Print("group (", i, ") with ", attacker.units, " attacking group (", target.initiative, ") with ", target.units, " with ", dmg)
				unitloss := dmg / target.hp
				byinitiative[target.initiative].units -= unitloss

				//normalise unitloss for pretty print & kill tracking =)
				if byinitiative[target.initiative].units < 0 {
					unitloss += byinitiative[target.initiative].units
				}
				//fmt.Println(" killing", unitloss)
				killed += unitloss
			}
		}

		if killed == 0 {
			//it's a draw, neither side is able to damage the other anymore.
			//#curveball
			return -1
		}
	}

	//add up winning army..
	count := 0
	for _, g := range r {
		if g.units > 0 {
			count += g.units
		}
	}

	return count
}

func main() {
	fmt.Println("Part1 : score ", part1(loaddata(testdata())))
	part2()
}

func testdata() string {
	return `
	Immune System:
	17 units each with 5390 hit points (weak to radiation, bludgeoning) with an attack that does 4507 fire damage at initiative 2
	989 units each with 1274 hit points (immune to fire; weak to bludgeoning, slashing) with an attack that does 25 slashing damage at initiative 3
	
	Infection:
	801 units each with 4706 hit points (weak to radiation) with an attack that does 116 bludgeoning damage at initiative 1
	4485 units each with 2961 hit points (immune to radiation; weak to fire, cold) with an attack that does 12 slashing damage at initiative 4	
	`
}

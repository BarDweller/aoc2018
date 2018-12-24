package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func loaddata2(input string, boost int) []group {
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
						g.attackpower = u + (boost * (1 - (state - 1))) //hacky hacky, side=1 == immune system, so 1-1-0 = 1 = boost, and side2=infection 1-2-1 = 0 = no boost
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

func immunewon(r []group) bool {
	for _, g := range r {
		if g.side == 1 && g.units > 0 {
			return true
		} else if g.side == 2 && g.units > 0 {
			return false
		}
	}
	panic("no winners")
}

func part2() {
	delta := 8192
	for i := delta; i > 1; {
		r := loaddata2(testdata(), i)
		score := part1(r)
		if score != -1 && immunewon(r) {
			if delta <= 1 {
				fmt.Println("Part 2: boost was", i, " score was ", score)
				break
			}
			delta = delta / 2
			i -= delta
		} else {
			//i was too small, go back to previous i
			i += delta
			delta = delta / 2
			i -= delta
		}
	}
}

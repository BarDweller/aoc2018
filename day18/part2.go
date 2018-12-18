package main

import "fmt"

func score(g grid) int {
	c := map[rune]int{}
	for _, v := range g.d {
		c[v]++
	}
	return (c['#'] * c['|'])
}

func checkmatch(scores [250]int, length int, idx1 int, idx2 int) bool {
	matchlen := idx2 - idx1
	//cant check if we don't have 2 full loops to compare with.
	if matchlen*2 > length {
		return false
	}

	//adjust indexes by array len, we'll always access using mod
	//so this becomes zero, but it keeps the array index > 0
	idx1 += len(scores)
	idx2 += len(scores)

	loops := (length / matchlen) - 1
	for loop := 0; loop < loops; loop++ {
		for check := 0; check < matchlen; check++ {
			base := (idx2 - check) % len(scores)
			test := (idx1 - (loop * matchlen) - check) % len(scores)
			if scores[base] != scores[test] {
				return false
			}
		}
	}
	//still here? we just verified a match for every full loop..
	//there may remain a partial loop we can test too..
	if length%matchlen != 0 {
		for check := 0; check < length%matchlen; check++ {
			base := (idx2 - check) % len(scores)
			test := (idx1 - (loops * matchlen) - check) % len(scores)
			if scores[base] != scores[test] {
				return false
			}
		}
	}
	//idx1 to idx2 identified a pattern that repeated throughout length elements of scores
	return true
}

func huntforpatterns(scores [250]int, count int) int {
	start := count % len(scores)
	length := count + 1
	if length > len(scores) {
		length = len(scores)
	}
	seek := scores[count%len(scores)]
	for i := 1; i < length; i++ {
		index := (len(scores) + start - i) % len(scores)
		if scores[index] == seek {
			//found two values that matched.. see if the values between repeat in the entire array
			found := checkmatch(scores, length, (len(scores) + start - i), len(scores)+start)
			if found {
				//if the values repeated, the array is full of a pattern of length i
				return i
			}
		}
	}
	return -1
}

func scoreforminute(scores [250]int, patternlen int, zeroidx int, zerominute int, minute int) int {
	minutessincezero := minute - zerominute
	//add delta to array index
	offset := zeroidx + (minutessincezero % patternlen)
	//wrap array index if requried
	offset = (offset) % len(scores)
	//return score for index
	return scores[offset]
}

func part2(g grid) {
	//buffer to hunt, tune size as required
	var scores [250]int
	//initial iterations to use to saturate buffer.
	max := 2000
	//fill buffer
	for i := 1; i < max; i++ {
		g = doMinute(g)
		scores[i%len(scores)] = score(g)
	}
	//hunt in buffer to see if a pattern has established.
	//last minute processed will be max-1
	matchlen := huntforpatterns(scores, max-1)
	if matchlen == -1 {
		panic("Did not find a pattern in the results, maybe try larger buffer, or more iters")
	}
	//verification phase =)
	//make sure our values per minute now match the real ones, go for a few more patten lengths to be sure.
	for i := max; i < max+(matchlen*3)+1; i++ {
		g = doMinute(g)
		actual := score(g)
		predicted := scoreforminute(scores, matchlen, (max-1-matchlen)%len(scores), max-1, i)
		if actual != predicted {
			panic("verification failed.")
		}
	}
	//return result for puzzle iter
	fmt.Println("Puzzle Iter Score ", scoreforminute(scores, matchlen, (max-1-matchlen)%len(scores), max-1, 1000000000))
}

func main2() {
	g := loaddata(testdata())
	part2(g)
}


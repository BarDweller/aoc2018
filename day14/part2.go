package main

import (
	"fmt"
)

func seekpattern(recipes []byte, currentrecipe []int, seek string) int {
	var seekbytes = make([]byte, len(seek))
	for i, r := range seek {
		seekbytes[i] = byte(r - '0')
	}
	seeklen := len(seek)
	for i := 0; true; i++ {
		oldcount := len(recipes)
		recipes = createnewrecipes(&recipes, &currentrecipe)
		newcount := len(recipes)
		if oldcount >= seeklen {
			delta := newcount - oldcount
			for d := 0; d < delta; d++ {
				found := true
				start := newcount - d
				for x, y := newcount-d-1, seeklen-1; y >= 0; x, y = x-1, y-1 {
					if seekbytes[y] != recipes[x] {
						found = false
						break
					}
				}
				if found {
					return start - seeklen
				}
			}
		}
	}
	return 0
}

func main2() {
	recipes := []byte{3, 7}
	currentrecipe := []int{0, 1}

	fmt.Println(seekpattern(recipes, currentrecipe, "59414"))
}

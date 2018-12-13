package main

import (
	"fmt"
	"sort"
)

func main2() {
	carts, tracks := loadgrid(data())

	for crashed(carts) < len(carts)-1 {
		sort.Slice(carts, func(i, j int) bool {
			switch {
			case carts[i].location.y < carts[j].location.y:
				return true
			case carts[i].location.y > carts[j].location.y:
				return false
			default:
				return carts[i].location.x < carts[j].location.x
			}
		})

		tick(carts, tracks)
	}

	for _, c := range carts {
		if !c.hasCrashed {
			fmt.Println(c)
		}
	}
}

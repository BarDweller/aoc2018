package main

import "fmt"

type pair struct {
	a, b byte
}

type recipes []byte

var transformmap map[pair]recipes

func buildtransforms() (result map[pair]recipes) {
	result = map[pair]recipes{}
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			source := pair{byte(i), byte(j)}
			total := byte(i + j)
			if total < 10 {
				result[source] = []byte{total}
			} else {
				result[source] = []byte{total / 10, total % 10}
			}
		}
	}
	return result
}

func createnewrecipes(recipes *[]byte, currentrecipes *[]int) []byte {
	current := pair{(*recipes)[(*currentrecipes)[0]], (*recipes)[(*currentrecipes)[1]]}
	sum := transformmap[current]
	result := append(*recipes, sum...)
	(*currentrecipes)[0] = ((*currentrecipes)[0] + 1 + int(current.a)) % len(result)
	(*currentrecipes)[1] = ((*currentrecipes)[1] + 1 + int(current.b)) % len(result)
	return result
}

func ppresult(recipes []byte, current []int) {
	for i, r := range recipes {
		switch {
		case i == current[0]:
			fmt.Print("(", r, ")")
		case i == current[1]:
			fmt.Print("[", r, "]")
		default:
			fmt.Print(" ", r, " ")
		}
	}
	fmt.Println("")
}

func seekten(recipes []byte, currentrecipe []int, seek int) string {
	for i := 0; len(recipes) < 10+seek; i++ {
		recipes = createnewrecipes(&recipes, &currentrecipe)
		//ppresult(recipes, currentrecipe)
	}
	result := ""
	for i := seek; i < seek+10; i++ {
		result += string('0' + recipes[i])
	}
	return result
}

func main() {
	transformmap = buildtransforms()

	recipes := []byte{3, 7}
	currentrecipe := []int{0, 1}
	fmt.Println(seekten(recipes, currentrecipe, 2018))

	main2()
}

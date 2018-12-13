package main

import (
	"fmt"
	"sort"
	"strings"
)

type xy struct {
	x, y int
}

type tracks struct {
	maxx int
	maxy int
	grid map[xy]rune
}

type direction int

const (
	up    direction = 0
	down  direction = 1
	left  direction = 2
	right direction = 3
)

type heading int

const (
	headedleft     heading = 0
	headedstraight heading = 1
	headedright    heading = 2
)

func (d direction) String() string {
	switch d {
	case up:
		return "^"
	case down:
		return "v"
	case left:
		return "<"
	case right:
		return ">"
	}
	panic("unknown enum")
}

func (h heading) String() string {
	switch h {
	case headedleft:
		return "L"
	case headedright:
		return "R"
	case headedstraight:
		return "S"
	}
	panic("unknown enum h")
}

type cart struct {
	dir         direction
	lastheading heading
	location    xy
	hasCrashed  bool
}

func (c cart) String() string {
	t := 'N'
	if c.hasCrashed {
		t = 'Y'
	}
	return fmt.Sprintf("{%s,%s,[%d,%d],%s} ", c.dir.String(), c.lastheading.String(), c.location.x, c.location.y, string(t))
}

func loadgrid(input string) (carts []cart, t tracks) {
	t = tracks{0, 0, map[xy]rune{}}
	x, y := 0, 0

	for _, line := range strings.Split(input, "\n") {
		//line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		for _, r := range line {
			switch r {
			case '>':
				{
					r = '-'
					carts = append(carts, cart{right, headedleft, xy{x, y}, false})
				}
			case '<':
				{
					r = '-'
					carts = append(carts, cart{left, headedleft, xy{x, y}, false})
				}
			case '^':
				{
					r = '|'
					carts = append(carts, cart{up, headedleft, xy{x, y}, false})
				}
			case 'v':
				{
					r = '|'
					carts = append(carts, cart{down, headedleft, xy{x, y}, false})
				}
			}
			if r != ' ' {
				t.grid[xy{x, y}] = r
			}
			x++
		}
		if x > t.maxx {
			t.maxx = x
		}
		y++
		x = 0
	}
	t.maxy = y - 1
	return carts, t
}

func pptracks(g tracks) {
	for y := 0; y <= g.maxy; y++ {
		for x := 0; x <= g.maxx; x++ {
			r, present := g.grid[xy{x, y}]
			if !present {
				r = ' '
			}
			fmt.Print(string(r))
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func pptracksandcarts(g tracks, c []cart) {
	locations := map[xy]*cart{}
	for i, cart := range c {
		locations[cart.location] = &c[i]
	}
	for y := 0; y <= g.maxy; y++ {
		for x := 0; x <= g.maxx; x++ {
			r, present := g.grid[xy{x, y}]
			if !present {
				r = ' '
			}
			if cart, hasCart := locations[xy{x, y}]; hasCart {
				r = rune(cart.dir.String()[0])
				if cart.hasCrashed {
					r = 'X'
				}
			}
			fmt.Print(string(r))
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func tick(carts []cart, t tracks) {
	locations := map[xy]*cart{}
	for i := range carts {
		if !carts[i].hasCrashed {
			locations[carts[i].location] = &carts[i]
		}
	}

	for i := range carts {
		if carts[i].hasCrashed {
			continue
		}
		switch carts[i].dir {
		case up:
			{
				next := xy{carts[i].location.x, carts[i].location.y - 1}
				delete(locations, carts[i].location)
				carts[i].location = next
				switch t.grid[next] {
				case '|':
					{
						//noop
					}
				case '\\':
					{
						carts[i].dir = left
					}
				case '/':
					{
						carts[i].dir = right
					}
				case '+':
					{
						switch carts[i].lastheading {
						case headedleft:
							{
								carts[i].dir = left
							}
						case headedright:
							{
								carts[i].dir = right
							}
						}
						carts[i].lastheading = (carts[i].lastheading + 1) % 3
					}
				default:
					{
						panic("Inconsistent map")
					}
				}

			}
		case down:
			{
				next := xy{carts[i].location.x, carts[i].location.y + 1}
				carts[i].location = next
				switch t.grid[next] {
				case '|':
					{
						//noop
					}
				case '\\':
					{
						carts[i].dir = right
					}
				case '/':
					{
						carts[i].dir = left
					}
				case '+':
					{
						switch carts[i].lastheading {
						case headedleft:
							{
								carts[i].dir = right
							}
						case headedright:
							{
								carts[i].dir = left
							}
						}
						carts[i].lastheading = (carts[i].lastheading + 1) % 3
					}
				default:
					{
						panic("Inconsistent map")
					}
				}
			}
		case left:
			{
				next := xy{carts[i].location.x - 1, carts[i].location.y}
				carts[i].location = next
				switch t.grid[next] {
				case '-':
					{
						//noop
					}
				case '\\':
					{
						carts[i].dir = up
					}
				case '/':
					{
						carts[i].dir = down
					}
				case '+':
					{
						switch carts[i].lastheading {
						case headedleft:
							{
								carts[i].dir = down
							}
						case headedright:
							{
								carts[i].dir = up
							}
						}
						carts[i].lastheading = (carts[i].lastheading + 1) % 3
					}
				default:
					{
						panic("Inconsistent map")
					}
				}
			}
		case right:
			{
				next := xy{carts[i].location.x + 1, carts[i].location.y}
				carts[i].location = next
				switch t.grid[next] {
				case '-':
					{
						//noop
					}
				case '\\':
					{
						carts[i].dir = down
					}
				case '/':
					{
						carts[i].dir = up
					}
				case '+':
					{
						switch carts[i].lastheading {
						case headedleft:
							{
								carts[i].dir = up
							}
						case headedright:
							{
								carts[i].dir = down
							}
						}
						carts[i].lastheading = (carts[i].lastheading + 1) % 3
					}
				default:
					{
						panic("Inconsistent map")
					}
				}
			}
		}
		if other, cartPresent := locations[carts[i].location]; cartPresent {
			//cart has crashed.
			other.hasCrashed = true
			carts[i].hasCrashed = true
			locations[carts[i].location] = &carts[i]
			fmt.Println("Crash! ", carts[i].String())
		} else {
			locations[carts[i].location] = &carts[i]
		}
	}
}

func crashed(carts []cart) int {
	count := 0
	for _, c := range carts {
		if c.hasCrashed {
			count++
		}
	}
	return count
}

func main() {
	carts, tracks := loadgrid(data())

	for crashed(carts) == 0 {
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
	fmt.Println("")

	main2()
}

func data() string {
	return `
/->-\        
|   |  /----\
| /-+--+-\  |
| | |  | v  |
\-+-/  \-+--/
  \------/   
	`
}

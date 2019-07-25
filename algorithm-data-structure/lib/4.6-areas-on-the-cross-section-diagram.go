package lib

import "fmt"

func ExecAreasOnTheCrossSectionDiagram() {
	args := `\\///\_/\/\\\\/_/\\///__\\\_\\/_\/_/\`
	runes := []rune(args)

	stack := make([]int, 0, len(runes)+1)
	top := 0

	push := func(x int) {
		top++
		stack[top] = x
	}

	pop := func() int {
		x := stack[top]
		top--
		return x
	}

	var area int
	for i, r := range runes {
		char := string(r)
		fmt.Println(i, char)
		switch char {
		case `\`:
			push(i)
		case `/`:
			x := pop()
			area += i - x
		default:
		}
	}
	fmt.Println(area)
}

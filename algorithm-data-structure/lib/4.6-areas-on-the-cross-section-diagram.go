package lib

import (
	"fmt"
	"sort"
)

func ExecAreasOnTheCrossSectionDiagram() {
	args := `\\///\_/\/\\\\/_/\\///__\\\_\\/_\/_/\`
	runes := []rune(args)

	routeStack := newRouteStack(len(runes))
	areaStack := newAreaStack(len(runes) / 2)

	var char string
	var areaSum int
	for cur, r := range runes {
		char = string(r)
		if char == `\` {
			routeStack.push(cur)
		} else if char == `/` && routeStack.top > 0 {
			left := routeStack.pop()
			areaSum += cur - left
			areaPart := cur - left
			for {
				if areaStack.top == 0 {
					break
				}
				if areaStack.data[areaStack.top].left <= left {
					break
				}
				old := areaStack.pop()
				areaPart += old.area
			}
			areaStack.push(pair{left, areaPart})
		}
	}
	fmt.Println(areaSum)
	fmt.Print(areaStack.top, " ")
	for i := 1; i <= areaStack.top; i++ {
		fmt.Print(areaStack.data[i].area, " ")
	}
	fmt.Print("\n")
}

type routeStack struct {
	top  int
	data []int
}

func newRouteStack(size int) *routeStack {
	return &routeStack{0, make([]int, size+1)}
}

func (s *routeStack) push(x int) {
	s.top++
	s.data[s.top] = x
}

func (s *routeStack) pop() int {
	x := s.data[s.top]
	s.top--
	return x
}

type pair struct {
	left int
	area int
}

type areaStack struct {
	top  int
	data []pair
}

func newAreaStack(size int) *areaStack {
	return &areaStack{0, make([]pair, size+1)}
}

func (s *areaStack) push(x pair) {
	s.top++
	s.data[s.top] = x
}

func (s *areaStack) pop() pair {
	x := s.data[s.top]
	s.top--
	return x
}

func ExecAreasOnTheCrossSectionDiagramOptimized() {
	args := `\\///\_/\/\\\\/_/\\///__\\\_\\/_\/_/\`
	AreasOnTheCrossSectionDiagramOptimized(args)
	AreasOnTheCrossSectionDiagramOptimized_2(args)
}

func AreasOnTheCrossSectionDiagramOptimized(args string) {
	runes := []rune(args)
	bars := make([]int, 0, len(runes))
	areas := map[int]int{}

	push := func(pos int) {
		bars = append(bars, pos)
	}

	pop := func() (int, bool) {
		n := len(bars)
		if n == 0 {
			return 0, false
		}

		x := bars[n-1]
		bars = bars[:n-1]
		return x, true
	}

	for i, r := range runes {
		char := string(r)
		switch char {
		case `\`:
			push(i)
		case `/`:
			if pos, ok := pop(); ok {
				area := i - pos
				for key, val := range areas {
					if pos < key && key < i {
						area += val
						delete(areas, key)
					}
				}
				areas[pos] = area
			}
		case `_`:
		}
	}

	keys := make([]int, 0, len(areas))
	for key, _ := range areas {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	var sum int
	for _, area := range areas {
		sum += area
	}

	fmt.Println(sum)
	fmt.Print(len(areas), " ")
	for _, key := range keys {
		fmt.Print(areas[key], " ")
	}
	fmt.Print("\n")
}

func AreasOnTheCrossSectionDiagramOptimized_2(args string) {
	runes := []rune(args)

	bars := make([]int, 0, len(runes))
	push := func(pos int) {
		bars = append(bars, pos)
	}
	pop := func() (int, bool) {
		n := len(bars)
		if n == 0 {
			return 0, false
		}

		x := bars[n-1]
		bars = bars[:n-1]
		return x, true
	}

	areas := make([][]int, 0, len(runes))
	pushA := func(x []int) {
		areas = append(areas, x)
	}
	popA := func() ([]int, bool) {
		n := len(areas)
		if n == 0 {
			return nil, false
		}

		x := areas[n-1]
		areas = areas[:n-1]
		return x, true
	}

	for i, r := range runes {
		char := string(r)
		switch char {
		case `\`:
			push(i)
		case `/`:
			if pos, ok := pop(); ok {
				area := i - pos
				for range areas {
					n := len(areas)
					if pos < areas[n-1][0] && areas[n-1][0] < i {
						under, _ := popA()
						area += under[1]
					}
				}
				pushA([]int{pos, area})
			}
		case `_`:
		}
	}

	var sum int
	for _, area := range areas {
		sum += area[1]
	}

	fmt.Println(sum)
	fmt.Print(len(areas), " ")
	for _, area := range areas {
		fmt.Print(area[1], " ")
	}
	fmt.Print("\n")
}

package lib

import "fmt"

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

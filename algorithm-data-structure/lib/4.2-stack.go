package lib

import (
	"fmt"
	"strconv"
)

type Stack struct {
	Top  int
	Data []int
}

func NewStack(size int) Stack {
	if size < 0 {
		size = 0
	}
	return Stack{
		Top:  0,
		Data: make([]int, size+1),
	}
}

func (s Stack) Initialize() {
	s.Top = 0
}

func (s Stack) IsEmpty() bool {
	return s.Top == 0
}

func (s Stack) IsFull() bool {
	return s.Top >= len(s.Data)-1
}

func (s Stack) Push(x int) (ok bool) {
	if s.IsFull() {
		return
	}
	s.Top++
	s.Data[s.Top] = x
	ok = true
	return
}

func (s Stack) Pop(x int) (ans int, ok bool) {
	if s.IsEmpty() {
		return
	}
	ans = s.Data[s.Top]
	ok = true
	s.Top--
	return
}

func ExecStack() {
	args := []string{"1", "2", "+", "3", "4", "-", "*"}
	top := 0
	S := make([]int, 1000)

	push := func(x int) {
		top++
		S[top] = x
	}

	pop := func() int {
		top--
		return S[top+1]
	}

	var a, b int

	for _, val := range args {
		switch val {
		case "+":
			b = pop()
			a = pop()
			fmt.Println(a, "+", b)
			push(a + b)
		case "-":
			b = pop()
			a = pop()
			fmt.Println(a, "-", b)
			push(a - b)
		case "*":
			b = pop()
			a = pop()
			fmt.Println(a, "*", b)
			push(a * b)
		default:
			i, err := strconv.Atoi(val)
			if err != nil {
				panic(err)
			}
			push(i)
		}
	}

	fmt.Println(pop())
}

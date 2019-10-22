package main

import "fmt"

func main() {
	l := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, n := range l {
		fmt.Println(EvenByModulo(n), EvenByBit(n))
	}
}

func EvenByModulo(n int) bool {
	return n%2 == 0
}

func EvenByBit(n int) bool {
	return (n & 1) == 0
}

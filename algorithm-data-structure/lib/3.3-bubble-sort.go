package lib

import "fmt"

func ExecBubbleSort() {
	fmt.Println(BubbleSort(5, 3, 2, 4, 1))
}

func BubbleSort(args ...int) ([]int, int) {
	swapped := 0
	finished := false
	for i := 0; !finished; i++ {
		finished = true
		for j := len(args) - 1; j > i; j-- {
			if args[j-1] > args[j] {
				args[j-1], args[j] = args[j], args[j-1]
				swapped++
				finished = false
			}
		}
	}
	return args, swapped
}

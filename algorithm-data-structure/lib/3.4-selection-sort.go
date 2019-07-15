package lib

import "fmt"

func ExecSelectionSort() {
	fmt.Println(SelectionSort(5, 6, 4, 2, 1, 3))
}

func SelectionSort(args ...int) ([]int, int) {
	swapped := 0
	for i := 0; i < len(args); i++ {
		min := i
		for j := i; j < len(args); j++ {
			if args[j] < args[min] {
				min = j
			}
		}
		if min != i {
			args[i], args[min] = args[min], args[i]
			swapped++
		}
	}
	return args, swapped
}

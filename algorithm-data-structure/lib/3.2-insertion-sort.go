package lib

import "fmt"

func ExecInsertionSort() {
	fmt.Println(InsertionSort(5, 2, 4, 6, 1, 3))
}

func InsertionSort(args ...int) []int {
	// fmt.Println(args)
	for i := 1; i < len(args); i++ {
		for j := i; j > 0; j-- {
			if args[j-1] > args[j] {
				args[j-1], args[j] = args[j], args[j-1]
			} else {
				break
			}
		}
		// fmt.Println(args)
	}
	return args
}

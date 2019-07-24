package lib

import "fmt"

func ExecPartition() {
	args := []int{13, 19, 9, 5, 12, 8, 7, 4, 21, 2, 6, 11}
	fmt.Println(partition(args, 0, 11))
}

func partition(args []int, p int, r int) ([]int, int) {
	x := args[r]
	i := p - 1
	for j := p; j < r; j++ {
		if args[j] <= x {
			i++
			args[i], args[j] = args[j], args[i]
		}
		fmt.Printf("i: %d, j: %d, args: %v\n", i, j, args)
	}
	args[i+1], args[r] = args[r], args[i+1]
	fmt.Printf("args: %v\n", args)
	return args, i + 1
}

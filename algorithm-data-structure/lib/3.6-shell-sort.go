package lib

import "fmt"

func ExecShellSort() {
	fmt.Println(ShellSort(5, 1, 4, 3, 2))
	fmt.Println(ShellSort(5, 1, 4, 3, 2, 4, 5, 6, 3, 4, 5, 4, 5, 5, 10, 11, 31, 2, 6, 7))
}

var cnt int

func InsertionSortWithG(g int, nums []int) []int {
	for i := g; i < len(nums); i++ {
		right := i
		left := right - g
		for left >= 0 && nums[left] > nums[right] {
			nums[left], nums[right] = nums[right], nums[left]
			left, right = left-g, left
			cnt++
		}
	}
	return nums
}

func ShellSort(nums ...int) []int {
	// 数列G = {1, 4, 13, 40, 121, 1093, ...} を生成
	G := make([]int, 0, len(nums))
	h := 1
	for h <= len(nums) {
		G = append(G, h)
		h = 3*h + 1
	}

	fmt.Print("G: ")
	for i := len(G) - 1; i >= 0; i-- {
		fmt.Print(G[i], " ")
		InsertionSortWithG(G[i], nums)
	}
	fmt.Println()
	fmt.Println("cnt:", cnt)
	return nums
}

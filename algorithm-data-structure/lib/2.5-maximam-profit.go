package lib

import (
	"fmt"
	"math"
)

func ExecMaximamProfit() {
	fmt.Println(MaximamProfit(5, 3, 1, 3, 4, 3))
	fmt.Println(MaximamProfit(4, 3, 2))

}

func MaximamProfit(nums ...int) int {
	max := math.MinInt32
	for i := 0; i < len(nums)-1; i++ {
		for j := i + 1; j < len(nums); j++ {
			max = int(math.Max(float64(max), float64(nums[j]-nums[i])))
		}
	}
	return max
}

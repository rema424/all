package lib

import "testing"

func TestMaximamProfit(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{"1", args{[]int{5, 3, 1, 3, 4, 3}}, 3},
		{"2", args{[]int{4, 3, 2}}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaximamProfit(tt.args.nums...); got != tt.want {
				t.Errorf("MaximamProfit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_maximamProfit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MaximamProfit(5, 3, 1, 3, 4, 3)
	}
}

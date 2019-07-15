package lib

import (
	"reflect"
	"testing"
)

func TestInsertionSort(t *testing.T) {
	type args struct {
		args []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
		{"1", args{[]int{5, 2, 4, 6, 1, 3}}, []int{1, 2, 3, 4, 5, 6}},
		{"2", args{[]int{4, 4, 6, 3, 8, 6, 7, 7, 8, 9, 4, 5, 2}}, []int{2, 3, 4, 4, 4, 5, 6, 6, 7, 7, 8, 8, 9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InsertionSort(tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertionSort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkInsertionSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InsertionSort(5, 2, 4, 6, 1, 3)
	}
}

package lib

import (
	"reflect"
	"testing"
)

func TestBubbleSort(t *testing.T) {
	type args struct {
		args []int
	}
	tests := []struct {
		name  string
		args  args
		want  []int
		want1 int
	}{
		// TODO: Add test cases.
		{"1", args{[]int{5, 2, 4, 6, 1, 3}}, []int{1, 2, 3, 4, 5, 6}, 9},
		{"2", args{[]int{4, 4, 6, 3, 8, 6, 7, 7, 8, 9, 4, 5, 2}}, []int{2, 3, 4, 4, 4, 5, 6, 6, 7, 7, 8, 8, 9}, 32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := BubbleSort(tt.args.args...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BubbleSort() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BubbleSort() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func BenchmarkBubbleSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BubbleSort(5, 2, 4, 6, 1, 3)
	}
}

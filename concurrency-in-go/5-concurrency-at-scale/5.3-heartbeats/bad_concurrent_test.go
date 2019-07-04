package main

import (
	"testing"
	"time"
)

func TestDoWork_GeneratesAllNumbers_Bad(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	_, results := DoWork(done, intSlice...)

	for i, expected := range intSlice {
		select {
		case r := <-results:
			if r != expected {
				t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
			}
		case <-time.After(1 * time.Second): // <1>
			// 壊れたゴルーチンがテストでデッドロックを起こしてしまわないのに十分と思われる時間が
			// 経過したあとでタイムアウトしています。
			t.Fatal("test timed out")
		}
	}
}

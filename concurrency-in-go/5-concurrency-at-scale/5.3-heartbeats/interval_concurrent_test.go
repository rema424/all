package main

import (
	"testing"
	"time"
)

func TestDoWork_GeneratesAllNumbers_Interval(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	const timeout = 2 * time.Second
	heartbeat, results := IntervalDoWork(done, timeout/2, intSlice...)

	<-heartbeat // <4>

	i := 0
	for {
		select {
		case r, ok := <-results:
			if !ok {
				return
			} else if expected := intSlice[i]; r != expected {
				t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
			}
			i++
		case <-heartbeat: // <5>
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}
	}
}

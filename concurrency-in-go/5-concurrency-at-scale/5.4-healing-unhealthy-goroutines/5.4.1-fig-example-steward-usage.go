package main

import "time"

func execExampleStewardUsage() {
	type startGoroutineFn1 func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) (heartbeat <-chan interface{}) // <1>

	newSteward := func(
		timeout time.Duration,
		startGoroutine startGoroutineFn1,
	) startGoroutineFn1 { // <2>
		return
	}
}

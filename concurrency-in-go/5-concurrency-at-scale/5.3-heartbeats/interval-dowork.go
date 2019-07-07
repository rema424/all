package main

import (
	"fmt"
	"time"
)

func execIntervalDowork() {
	// done := make(chan interface{})
	// defer close(done)

	// heartbeat, results := IntervalDoWork(done, 40*time.Microsecond, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	// // for-select
	// for {
	// 	select {
	// 	case _, ok := <-heartbeat:
	// 		if ok {
	// 			fmt.Println("pulse")
	// 		} else {
	// 			return
	// 		}
	// 	case r, ok := <-results:
	// 		if ok {
	// 			fmt.Printf("results %v\n", r)
	// 		} else {
	// 			return
	// 		}
	// 	}
	// }

	done := make(chan interface{})
	defer close(done)

	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start).Seconds())
	}()

	intSlice := []int{0, 1, 2, 3, 5}
	const timeout = 10 * time.Second
	heartbeat, results := IntervalDoWork(done, timeout/2, intSlice...)

	// まだ最初のハートビートがゴルーチンのループに入ったことを知らせてくれるのを待っています。
	<-heartbeat // <4>

	// i := 0
	for {
		select {
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
			// i++
		case <-heartbeat: // <5>
			// タイムアウトが発生しないようにハートビートからもselectで値を取得しています。
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
}

func IntervalDoWork(
	done <-chan interface{},
	pulseInterval time.Duration,
	nums ...int,
) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)

	go func() {
		defer close(heartbeat)
		defer close(intStream)

		time.Sleep(2 * time.Second)

		pulse := time.Tick(pulseInterval)

		// ラベルを使って内側のループから少し簡単にcontinueできるようにしました。
	numLoop: // <2>
		for _, n := range nums {
			// 2つのループが必要です。1つは数のリストをrangeで回すためで、このコメントがある内側のループは
			// 数がintStreamに無事送信されるまで繰り返すものです。
			for { // <1>
				select {
				case <-done:
					return
				case <-pulse:
					select {
					case heartbeat <- struct{}{}:
						fmt.Println("push heartbeat", n)
					default:
					}
				case intStream <- n:
					fmt.Println("push", n)
					// 外側のループをcontinueしています。
					continue numLoop // <3>
				}
			}
		}
	}()

	return heartbeat, intStream
}

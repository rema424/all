package main

import "time"

func DoWork(done <-chan interface{}, nums ...int) (<-chan interface{}, <-chan int) {
	// 戻り値の宣言
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)

	go func() {
		defer close(heartbeat)
		defer close(intStream)

		// ここでゴルーチンを起動できる状態にする前に何らかの遅延をシミュレートしています。
		// 実際にはこの遅延はあらゆる理由で発生しうるもので、非決定的です。
		// 私が経験したものでは、CPU 負荷、ディスク競合、ネットワーク遅延、あとは妖精さんによるものがありました。
		time.Sleep(2 * time.Second) // <1>

		// for-select
		for _, n := range nums {
			select {
			case heartbeat <- struct{}{}:
			default:
				// バッファ対策 default
			}

			select {
			case <-done:
				return
			case intStream <- n:
			}
		}
	}()

	return heartbeat, intStream
}

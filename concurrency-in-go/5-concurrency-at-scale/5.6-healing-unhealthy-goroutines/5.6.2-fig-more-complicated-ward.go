package main

import "time"

func execMoreComplicatedWard() {

	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)

			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	bridge := func(
		done <-chan interface{},
		chanStream <-chan <-chan interface{},
	) <-chan interface{} {
		valStream := make(chan interface{}) // <1>
		go func() {
			defer close(valStream)
			for { // <2>
				var stream <-chan interface{}

				select {
				case maybeStream, ok := <-chanStream:
					if !ok {
						return
					}
					stream = maybeStream
				case <-done:
					return
				}

				for val := range orDone(done, stream) { // <3>
					select {
					case valStream <- val:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	type startGoroutineFn func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) (heartbeat <-chan interface{}) // <1>

	// 中庭に囲んでもらいたい値を取って、中庭からのやり取りに使うチャネルを返します。
	doWorkFn := func(
		done <-chan interface{},
		intList ...int,
	) (startGoroutineFn, <-chan interface{}) { // <1>

		// この行でブリッジパターンの一部としてチャネルのチャネルを作ります。
		intChanStream := make(chan (<-chan interface{})) // <2>

		// doWorkFnの第二戻り値
		intStream := bridge(done, intChanStream)

		// doWorkFnの第一戻り値
		// 管理人に監視されるクロージャーを作ります。
		doWork := func(
			done <-chan interface{},
			pulseInterval time.Duration,
		) <-chan interface{} { //<3>

			// ここで中庭のゴルーチンのインスタンスの中でやり取りするチャネルを初期化します。
			intStream := make(chan interface{}) // <4>
			heartbeat := make(chan interface{})
			go func() {
				defer close(intStream)

				select {
				case intChanStream <- intStream: // <5>
				case <-done:
					return
				}
			}()
			return heartbeat
		}

		return doWork, intStream
	}
}

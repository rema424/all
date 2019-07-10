package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func execMoreComplicatedWard() {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	var orDone func(done <-chan interface{}, c <-chan interface{}) <-chan interface{}
	var bridge func(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{}
	var take func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{}

	or = func(channels ...<-chan interface{}) <-chan interface{} { // <1>
		switch len(channels) {
		case 0: // <2>
			// 引数のdoneチャネルが0件の場合はnil
			return nil
		case 1: // <3>
			// 引数のdoneチャネルが1件の場合はそれを返却
			return channels[0]
		}

		// 引数のdoneチャネルが2件以上の場合は一旦ここで宣言するorDoneチャネルを返却する
		// goroutineの中で引数で渡ってきたdoneチャネルを監視して、どれか一つでもcloseされたら
		// 先に返却しておいたorDoneチャネルにcloseを通知する。
		orDone := make(chan interface{})
		go func() { // <4>
			defer close(orDone)

			switch len(channels) {
			case 2: // <5>
				// 引数のdoneチャネルが2件の場合は、そのどちらかに値が届いた時点で
				// selectを抜ける -> switchを抜ける -> 処理の終点に到達してdefer close()が呼ばれる -> goroutineを抜ける
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default: // <6>
				remainingCh := append(channels[3:], orDone)
				recursionCh := or(remainingCh...)
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-recursionCh: // <6>
				}
			}
		}()
		return orDone
	}

	orDone = func(done, c <-chan interface{}) <-chan interface{} {
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

	bridge = func(
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

	take = func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {

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

	newSteward := func(timeout time.Duration, startGoroutine startGoroutineFn) startGoroutineFn { // <2>
		return func(done <-chan interface{}, pulseInterval time.Duration) <-chan interface{} {
			heartbeat := make(chan interface{}) // return するやつを宣言
			go func() {
				defer close(heartbeat) // return するやつを close

				var wardDone chan interface{}        // newSteward の引数で渡ってきた startGoroutineFn に引数で渡すやつ
				var wardHeartbeat <-chan interface{} // newSteward の引数で渡ってきた startGoroutineFn の戻り値で返却されるやつ

				startWard := func() { // <3>
					wardDone = make(chan interface{})                             // <4>
					wardHeartbeat = startGoroutine(or(wardDone, done), timeout/2) // <5>
				}

				startWard()
				pulse := time.Tick(pulseInterval)

			monitorLoop:
				for {
					timeoutSignal := time.After(timeout)

					for { // <6>
						select {
						case <-pulse:
							select {
							case heartbeat <- struct{}{}:
							default:
							}
						case <-wardHeartbeat: // <7>
							continue monitorLoop
						case <-timeoutSignal: // <8>
							log.Println("steward: ward unhealthy; restarting")
							close(wardDone)
							startWard()
							continue monitorLoop
						case <-done:
							return
						}
					}
				}
			}()
			return heartbeat // return するやつ
		}
	}

	log.SetFlags(log.Ltime | log.LUTC)
	log.SetOutput(os.Stdout)

	done := make(chan interface{})
	defer close(done)

	doWork, intStream := doWorkFn(done, 1, 2, -1, 3, 4, 5)      // <1>
	doWorkWithSteward := newSteward(1*time.Millisecond, doWork) // <2>
	doWorkWithSteward(done, 1*time.Hour)                        // <3>

	for intVal := range take(done, intStream, 6) { // <4>
		fmt.Printf("Received: %v\n", intVal)
	}
}

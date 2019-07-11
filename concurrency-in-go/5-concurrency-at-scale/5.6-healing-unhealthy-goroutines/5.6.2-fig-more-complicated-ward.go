package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rema424/all/lib/concurrency"
)

func execMoreComplicatedWard() {
	type startGoroutineFn func(done <-chan interface{}, pulseInterval time.Duration) (heartbeat <-chan interface{}) // <1>

	// 中庭に囲んでもらいたい値を取って、中庭からのやり取りに使うチャネルを返します。
	doWorkFn := func(done <-chan interface{}, intList ...int) (doWork startGoroutineFn, intStream <-chan interface{}) { // <1>

		// この行でブリッジパターンの一部としてチャネルのチャネルを作ります。
		intChanStream := make(chan (<-chan interface{})) // <2>

		// doWorkFnの第二戻り値
		intStream = concurrency.Bridge(done, intChanStream)

		// doWorkFnの第一戻り値
		// 管理人に監視されるクロージャーを作ります。
		doWork = func(done <-chan interface{}, pulseInterval time.Duration) <-chan interface{} { //<3>
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

				pulse := time.Tick(pulseInterval)

				for {
				valueLoop:
					for _, intVal := range intList {
						if intVal < 0 {
							log.Printf("negative value: %v\n", intVal) // <6>
							return
						}

						for {
							select {
							case <-pulse:
								select {
								case heartbeat <- struct{}{}: // heartbeatに空きが出るまでブロック
								default: // heartbeatにすでに値が入っていたら何もしないでselectを抜ける
								}
							case intStream <- intVal:
								continue valueLoop
							case <-done:
								return
							}
						}
					}
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
					wardDone = make(chan interface{})                                         // <4>
					wardHeartbeat = startGoroutine(concurrency.Or(wardDone, done), timeout/2) // <5>
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

	// この行では中庭の関数を作ります。
	doWork, intStream := doWorkFn(done, 1, 2, 3, 4, 5, -1, 7) // <1>
	// ここでdoWorkのクロージャーを監視する管理人を作成します。失敗がかなりすぐ発生することが予想されるので、監視機関を1ミリ秒だけ設定します。
	// doWorkWithSteward := newSteward(1*time.Millisecond, doWork) // <2>
	// 3秒たったらリスタートする。
	doWorkWithSteward := newSteward(2*time.Second, doWork) // <2>
	// 管理人に中庭を起動して監視を開始するように伝えます。ハートビートが1時間です。
	// 戻り値を変数に代入していませんが、この関数はheatbeatを返却しています。
	doWorkWithSteward(done, 1*time.Hour) // <3>

	// 6個数値が流れてくるまで処理を続ける。
	// 流れてくるまで管理人がdoWorkをリスタートしてくれる。
	for intVal := range concurrency.Take(done, intStream, 20) { // <4>
		fmt.Printf("Received: %v\n", intVal)
	}
}

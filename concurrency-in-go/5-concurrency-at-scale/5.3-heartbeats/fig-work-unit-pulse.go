package main

import (
	"fmt"
	"math/rand"
)

func execWorkUnitPulse() {
	// -----------------
	// doWork 関数 start
	// -----------------
	doWork := func(done <-chan interface{}) (<-chan interface{}, <-chan int) {
		// 戻り値を宣言

		// ここでバッファが 1 の heartbeat チャネルを作成します。
		// これで送信待ちのものが何もなくても最低 1 つの鼓動が常に送られることを保証しています。
		heartbeatStream := make(chan interface{}, 1) // <1>
		workStream := make(chan int)

		// 並行処理部分
		// 宣言済みのチャネルに値を流し込む
		go func() {
			// done と同時に colse する
			defer close(heartbeatStream)
			defer close(workStream)

			// for-select
			for i := 0; i < 10; i++ {
				// 特に待ちは発生しない
				// ここでハートビート用に異なる select ブロックを用意しています。
				// これを result への送信と同じ select ブロックに入れたくなかった理由は、
				// 受信側が結果を受け取る準備ができていなかったときはかわりにハートビートの鼓動を
				// 受信してしまい現在の結果の値を受け取り損ねてしまうtからです。
				// また done チャネルの条件も含めていません。
				// なぜならフォールスルーする default の条件があるからです。
				select { // <2>
				case heartbeatStream <- struct{}{}:
				default: // <3>
					// 再度何もハートビートを待っていない可能性に対処しています。
					// heartbeat チャネルはバッファを 1 で作っているため、もし何かがハートビートを
					// 待っているけれど、最初の鼓動には間に合わなかった場合、それでもなお鼓動の通知を受けます。
				}

				// 特に待ちは発生しない
				select {
				case <-done:
					return
				case workStream <- rand.Intn(10):
				}
			}
			// for 文の中で待ちは発生しないのですぐ終了する。
			// 終了時にチャネルが close される。
		}()
		// 宣言済みの戻り値を返却
		return heartbeatStream, workStream
	}
	// -----------------
	// doWork 関数 end
	// -----------------

	// ---------------
	// メイン処理 start
	// ---------------
	done := make(chan interface{})
	defer close(done)

	heartbeat, results := doWork(done)

	// for-select
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
}

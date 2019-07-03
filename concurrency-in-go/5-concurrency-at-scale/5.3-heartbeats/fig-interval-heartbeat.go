package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) (<-chan interface{}, <-chan time.Time) {
		// ハートビートを送信する先のチャネルを設定します。
		// このチャネルを doWork より返します。
		heartbeat := make(chan interface{}) // <1>
		results := make(chan time.Time)

		go func() {
			defer close(heartbeat)
			defer close(results)

			// doWork の引数で与えられた pulseInterval の周期でハートビートの鼓動を設定します。
			// pulseInterval ごとにこのチャネルでは何かしら読み込めるようにします。
			pulse := time.Tick(pulseInterval) // <2>
			// 仕事が入ってくる様子のシミュレートに使われる別のティッカーです。
			// ここでは pulseInterval よりも大きな周期を選びました。
			// これによりゴルーチンからハートビートにやってくるのを確認できます。
			workGen := time.Tick(2 * pulseInterval) // <3>

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default: // <4>
					// default 節を含めていることに注意してください。
					// 誰もハートビートを確認していない可能性があるということに対して
					// 常に対策をしなければなりません。
					// ゴルーチンから創出される結果は極めて重要ですが、ハートビートの鼓動はそうではありません。
				}
			}

			sendResult := func(r time.Time) {
				for {
					select {
					case <-done:
						// 親によりキャンセルされた場合
						return
					case <-pulse: // <5>
						// done チャネルと同様に、送信や受信を行うときはいつでもハートビートの
						// 鼓動に対する条件を含める必要があります。
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			// メインのブロッキング処理
			for {
				select {
				// キャンセル時
				case <-done:
					// 親によりキャンセルされた場合
					return
				// 1秒ごとに pulse に流れ込む
				case <-pulse: // <5>
					sendPulse()
				// 2秒ごとに workGen に流れ込む
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	}

	done := make(chan interface{})
	// 標準的な done チャネルを作り、10 秒後に閉じます。
	// これでゴルーチンになんらかの仕事を与えます。
	time.AfterFunc(10*time.Second, func() { close(done) }) // <1>

	// タイムアウト時間を設定します。
	// ハートビートの周期とタイムアウトを紐付けるためにこれを使います。
	const timeout = 2 * time.Second // <2>

	// ここで timeout/2 を渡します。これによりハートビートが追加で鼓動を打たせて、
	// タイムアウトに対して過敏にならなくてすむようにします。
	heartbeat, results := doWork(done, timeout/2) // <3>

	for {
		select {
		case _, ok := <-heartbeat: // <4>
			// ハートビートに対して select をかけます。結果が何もなかった場合、少なくとも
			// heartbeat チャネルから timeout/2 経過するごとにメッセージを受け取れることが保証されています。
			// もしハートビートからの鼓動が受け取れない場合は、ゴルーチン自体に何かしらの問題があるとわかります。
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results: // <5>
			// results チャネルより select をかけます。ここでは特に面白いものはありません。
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout): // <6>
			// ハートビートも新しい結果も受け取らなかった場合にはタイムアウトします。
			return
		}
	}
}

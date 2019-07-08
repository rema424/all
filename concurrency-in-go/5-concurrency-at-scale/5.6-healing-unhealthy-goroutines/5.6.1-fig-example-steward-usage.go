package main

import (
	"log"
	"os"
	"time"
)

func execExampleStewardUsage() {
	var or func(channels ...<-chan interface{}) <-chan interface{}

	or = func(channels ...<-chan interface{}) <-chan interface{} { // <1>
		switch len(channels) {
		case 0: // <2>
			return nil
		case 1: // <3>
			return channels[0]
		}

		orDone := make(chan interface{})

		go func() { // <4>
			defer close(orDone)

			switch len(channels) {
			case 2: // <5>
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default: // <6>
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...): // <6>
				}
			}
		}()

		return orDone
	}

	// 監視と再起動ができるゴルーチンのシグネチャを定義します。
	// 見慣れたdoneチャネルとハートビートパターンのpulseIntervalとheartbeatがあります。
	type startGoroutineFn1 func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) (heartbeat <-chan interface{}) // <1>

	// この行で管理人が監視するゴルーチンのためのtimeoutと、
	// 監視するゴルーチンを起動するためのstartGoroutine関数を取っています。
	// 興味深いことに、管理人自身はstartGoroutineFnを返していて、
	// これは管理人自身も監視可能であることそ示しています。
	newSteward := func(
		timeout time.Duration,
		startGoroutine startGoroutineFn1,
	) startGoroutineFn1 { // <2>

		// おなじみのハートビートパターン関数を返却する。
		// この関数の中で、引数で与えられたgoroutineで起動したい関数を起動する。
		return func(
			done <-chan interface{}, // 管理人自身のdoneチャネル
			pulseInterval time.Duration, // 管理人自身の鼓動期間
		) <-chan interface{} {
			heartbeat := make(chan interface{}) // 管理人自身のハートビート

			// この中で引数で与えられた関数を実行する
			go func() {
				// （引数で与えられた）実行したい関数が完了したら管理人のハートビートを閉じる
				defer close(heartbeat)

				// 中庭（引数で与えられた実行したい関数）のdoneチャネルとハートビート
				var wardDone chan interface{}
				var wardHeartbeat <-chan interface{}

				// 監視しているゴルーチンを起動するための一貫した方法としてクロージャーを定義しています。
				startWard := func() { // <3>
					// 停止すべきだとシグナルを送る必要がある場合に備えて、中庭のゴルーチンに渡す新しいチャネルを作成しています。
					wardDone = make(chan interface{}) // <4>
					// 監視対象のゴルーチンを起動します。
					// 管理人が停止するか、管理人が中庭のゴルーチンを停止させたい場合に対象のゴルーチンには停止してもらいたいので、
					// 両方のdoneチャネルをorの中に内包します。
					// 渡しているpulseIntervalはタイムアウト期間の半分です。
					// これは5.3ハートビートの説で話した様に調整可能ではあります。
					wardHeartbeat = startGoroutine(or(wardDone, done), timeout/2) // <5>
				}

				startWard() // 中庭が起動する
				pulse := time.Tick(pulseInterval)

			monitorLoop:
				for {
					// 中庭のハートビートはtimeout/2
					timeoutSignal := time.After(timeout)

					// 内側のループです。これは管理人が自身の鼓動を確実に外へと送信できるようにしています。
					for { // <6>
						select {
						case <-pulse:
							select {
							case heartbeat <- struct{}{}:
							default:
							}
						case <-wardHeartbeat: // <7>
							// 中庭の鼓動を受信したら、監視のループを継続する、という実装になっているのがわかります。
							continue monitorLoop
						case <-timeoutSignal: // <8>
							// タイムアウト期間内に中庭からの鼓動が受信できなければ。中庭に停止するようリクエストし、
							// 新しい中庭のゴルーチンを起動し始めることを示している行です。
							// その後、監視を続けます。
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

			return heartbeat
		}
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	doWork := func(done <-chan interface{}, _ time.Duration) <-chan interface{} {
		log.Println("ward: Helo, I'm irresponsible!")

		go func() {
			<-done // <1>
			log.Println("ward: I am halting.")
		}()

		return nil
	}

	doWorkWithSteward := newSteward(4*time.Second, doWork) // <2>

	done := make(chan interface{})
	time.AfterFunc(9*time.Second, func() { // <3>
		log.Println("main: halting steward and ward.")
		close(done)
	})

	for range doWorkWithSteward(done, 4*time.Second) {
	} // <4>

	log.Println("Done")
}

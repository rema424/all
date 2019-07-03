package main

func execWorkUnitPulse() {
	// -----------
	// doWork 関数
	// -----------
	doWork := func(done <-chan interface{}) (<-chan interface{}, <-chan int) {
		// 戻り値を宣言
		heartbeatStream := make(chan interface{})
		workStream := make(chan int)

		// 並行処理部分
		// 宣言済みのチャネルに値を流し込む
		go func() {
			// done と同時に colse する
			defer close(heartbeatStream)
			defer close(workStream)

			// for-select
			for i := 0; i < 10; i++ {
				select {}
			}
		}()

		// 宣言済みの戻り値を返却
		return heartbeatStream, workStream
	}
}

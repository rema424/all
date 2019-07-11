package concurrency

// Orは複数のdoneチャネルを引数にとります。
// このどれか一つでもcloseされたらcloseされるdoneチャネルを作成して返却します。
func Or(channels ...<-chan interface{}) <-chan interface{} { // <1>
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
			recursionCh := Or(remainingCh...)
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

// OrDoneはdoneチャネルと値チャネルを引数にとります。
// doneチャネルまたは値チャネルのどちらかが閉じられるまで値チャネルから値を読み込むチャネルを返却するヘルパー関数です。
func OrDone(done, c <-chan interface{}) <-chan interface{} {
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

// Bridgeはチャネルが流れてくるチャネルを扱いやすくするヘルパー関数です。
func Bridge(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
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

			for val := range OrDone(done, stream) { // <3>
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

// Take は valueStream から num 個の値を取り出すチャネルを返すヘルパー関数です。
func Take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

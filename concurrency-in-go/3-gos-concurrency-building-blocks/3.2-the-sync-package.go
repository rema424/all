package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

func execWaitGroup() {
	var wg sync.WaitGroup

	// Addを引数に1を渡して呼び出し、1つのゴルーチンが起動したことを表しています。
	wg.Add(1) // <1>
	go func() {
		// Doneをdeferキーワードを使って呼び出して、ゴルーチンのクロージャーが終了する前に
		// WaitGroupに終了することを確実に伝えるようにしています。
		defer wg.Done() // <2>
		fmt.Println("1st goroutine sleeping...")
		time.Sleep(1)
	}()

	wg.Add(1) // <1>
	go func() {
		defer wg.Done() // <2>
		fmt.Println("2nd goroutine sleeping...")
		time.Sleep(2)
	}()

	// Waitを呼び出しています。
	// 全てのゴルーチンが終了したと伝えるまでメインゴルーチンをブロックします。
	wg.Wait() // <3>
	fmt.Println("All goroutines complete.")
}

func execWaitGroup_2() {
	hello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Printf("Helo from %v!\n", id)
	}

	const numGreeters = 5
	var wg sync.WaitGroup
	wg.Add(numGreeters)
	for i := 0; i < numGreeters; i++ {
		go hello(&wg, i+1)
	}
	wg.Wait()
}

func execMutexAndRWMutex() {
	var count int
	var lock sync.Mutex

	increment := func() {
		lock.Lock()         // <1>
		defer lock.Unlock() // <2>
		count++
		fmt.Printf("Incrementing: %d\n", count)
	}

	decrement := func() {
		lock.Lock()         // <1>
		defer lock.Unlock() // <2>
		count--
		fmt.Printf("Decrementing: %d\n", count)
	}

	// インクリメント
	var arithmetic sync.WaitGroup
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			increment()
		}()
	}

	// デクリメント
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			decrement()
		}()
	}

	arithmetic.Wait()
	fmt.Println("Arithmetic complete.")
}

func execMutexAndRWMutex_2() {
	// producer関数の2番目の引数はwync.Locker型です。
	// このインタフェースにはLockとUnlockという2つのメソッドがあります。
	// このインタフェースはMutex型とRWMutex型を満たします。
	producer := func(wg *sync.WaitGroup, l sync.Locker) { // <1>
		defer wg.Done()
		for i := 5; i > 0; i-- {
			l.Lock()
			l.Unlock()
			// 生産者を1秒スリープさせて、Observerゴルーチンよりも非活発にします。
			time.Sleep(1) // <2>
		}
	}

	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}

	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1)
		beginTestTime := time.Now()
		go producer(&wg, mutex)
		for i := count; i > 0; i-- {
			go observer(&wg, rwMutex)
		}

		wg.Wait()
		return time.Since(beginTestTime)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()

	var m sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutex\tMutex\n")
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			test(count, &m, m.RLocker()),
			test(count, &m, &m),
		)
	}
}

func execCond() {
	conditionTrue := func() bool {
		return false
	}
	// 新しいCondのインスタンスを作ります。
	// NewCond関数はsync.Lockerインスタンスを満たす型を引数に取ります。
	// これによって、Cond型が他のゴルーチンを並行処理で安全な方法で協調できるようになります。
	c := sync.NewCond(&sync.Mutex{}) // <1>
	// この条件でLocekrをロックします。
	// これはWaitへの呼び出しがループに入るときに自動的にUnlockを呼び出すため必要です。
	c.L.Lock() // <2>
	for !conditionTrue() {
		// 条件が発生したかどうかが通知されるのを待ちます。これはブロックする呼び出しで、ゴルーチンは一時停止します。
		c.Wait() // <3>
	}
	// この条件でLockerのロックを解除します。この記述はWaitの呼び出しが終わると、
	// このじょうけんでLockを呼び出すので、必要です。
	c.L.Unlock() // <4>
}

func execCond_2() {
	// まず標準のsync.MutexをLockerとして使って条件を作成します。
	c := sync.NewCond(&sync.Mutex{}) // <1>
	// 次に、長さ0のスライスを作成します。最終的に10買い足すと分かっているので、キャパシティを10に設定します。
	queue := make([]interface{}, 0, 10) // <2>

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		// 再度条件のクリティカルセクションに入って、条件に合った形でデータを修正します。
		c.L.Lock() // <8>
		// スライスの先頭をスライスの2番目の要素を指すように帰ることでキューから取り出したことにします。
		queue = queue[1:] // <9>
		fmt.Println("Removed from queue")
		// 無事に要素をキューから取り出したので条件のクリティカルセクションを抜けます。
		c.L.Unlock() // <10>
		// 条件を待っているゴルーチンに何かが起きたことを知らせます。
		c.Signal() // <11>
	}

	for i := 0; i < 10; i++ {
		// 条件であるLockerのLockメソッドを呼び出してクリティカルセクションに入ります。
		c.L.Lock() // <3>
		// ループ内でキューの長さを確認します。これは重要なことです。
		// なぜなら条件上のシグナルは必ずしも同じ待っている事象が起きたことを意味していないからです
		// ーー何かが起きただけです。
		for len(queue) == 2 { // <4>
			// Waitを呼び出します。これによって条件のシグナルが創出されるまでメインゴルーチンを一時停止します。
			c.Wait() // <5>
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		// 1秒後に要素をキューから取り出す新しいゴルーチンを生成します。
		go removeFromQueue(1 * time.Second) // <6>
		// 要素を無事キューに追加できたので条件のクリティカルセクションを抜けます。
		c.L.Unlock() // <7>
	}
}

func execCond_3() {
	// Clickedという条件を含んでいるButton型を定義します。
	type Button struct { // <1>
		Clicked *sync.Cond
	}
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	// 条件に応じて送られてくるシグナルを扱う関数を登録するための便利な関数を定義します。
	// 各ハンドラーはそれぞれのゴルーチン上で動作します。
	// そしてsubscribeはゴルーチンが実行されていると確認できるまで終了しません。
	subscribe := func(c *sync.Cond, fn func()) { // <2>
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}

	// WaitGroupを作ります。これはプログラムがstdoutへ書き込む前に終了してしまわないようにするためだけのものです。
	var clickRegisterd sync.WaitGroup // <3>
	clickRegisterd.Add(3)

	// マウスのボタンが離された時のハンドラーを設定します。
	// こちらはClickedという状態（Cond）に対応するBroadcastを呼び出して、全てのハンドラーにマウスのボタンがクリックしたということを知らせます。
	// （より堅牢な実装では最初にボタンが押下されたかを確認すれば良いでしょう。）

	subscribe(button.Clicked, func() { // <4>
		fmt.Println("Maximizing window.")
		clickRegisterd.Done()
	})

	subscribe(button.Clicked, func() { // <5>
		fmt.Println("Displaying annoying dialog box!")
		clickRegisterd.Done()
	})

	subscribe(button.Clicked, func() { // <6>
		fmt.Println("Mouse clicked.")
		clickRegisterd.Done()
	})

	time.Sleep(2 * time.Second)
	button.Clicked.Broadcast() // <7>
	clickRegisterd.Wait()
}

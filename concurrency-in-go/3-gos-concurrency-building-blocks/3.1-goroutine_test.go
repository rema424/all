package main

import (
	"sync"
	"testing"
)

func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	var token struct{}
	sender := func() {
		defer wg.Done()
		// 開始と言われるまで待機します。
		// コンテキストスイッチの計測にゴルーチンの設定と起動の時間を入れたくなかったのでこうしています。
		<-begin // <1>
		for i := 0; i < b.N; i++ {
			// 受信側のゴルーチンにメッセージを送信しています。
			// struct{}はから構造体と呼ばれるもので、メモリを消費しません。
			// これによって、メッセージを送出する時間だけを計測できます。
			c <- token // <2>
		}
	}
	receiver := func() {
		defer wg.Done()
		<-begin // <1>
		for i := 0; i < b.N; i++ {
			// メッセージを受信しますが、何もしません。
			<-c // <3>
		}
	}

	wg.Add(2)
	go sender()
	go receiver()
	// タイマーを起動します。
	b.StartTimer() // <4>
	// 2つのゴルーチンに開始を伝えます。
	close(begin) // <5>
	wg.Wait()
}

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func execDuplicatedRequest() {
	done := make(chan interface{})
	result := make(chan int)

	var wg sync.WaitGroup
	wg.Add(10)

	// ここでリクエストを扱う10個のハンドラーを起動します。
	for i := 0; i < 10; i++ { // <1>
		go duplicatedDoWork(done, i, &wg, result)
	}

	// この業ではハンドラー群の中から最初に返された値を取り出します。
	firstReturned := <-result // <2>

	// 残り全てのハンドラーをキャンセルします。こうすることで不必要な仕事をし続けないようにします。
	close(done) // <3>
	wg.Wait()

	fmt.Printf("Received an answer from #%v\n", firstReturned)
}

func duplicatedDoWork(
	done <-chan interface{},
	id int,
	wg *sync.WaitGroup,
	result chan<- int,
) {
	started := time.Now()
	defer wg.Done()

	// Simulate random load
	simulatedLoadTime := time.Duration(1+rand.Intn(5)) * time.Second

	// ブロックしているだけでdoneの後に処理を中断(return)しているわけではない。
	// 処理を続けて、処理が完了されたら defer wg.Done() が呼ばれる。
	select {
	case <-done:
	case <-time.After(simulatedLoadTime):
	}

	// ブロックしているだけでdoneの後に処理を中断(return)しているわけではない。
	// 処理を続けて、処理が完了されたら defer wg.Done() が呼ばれる。
	select {
	case <-done:
	case result <- id:
	}

	took := time.Since(started)
	// Display how long handlers would have taken
	if took < simulatedLoadTime {
		took = simulatedLoadTime
	}
	fmt.Printf("%v took %v\n", id, took)
}

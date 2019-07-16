package main

import (
	"fmt"
	"runtime"
	"sync"
)

func execGoroutine() {
	go sayHello()
	// 他の処理を続ける
}

func sayHello() {
	fmt.Println("hello")
}

func execGoroutine_2() {
	go func() {
		fmt.Println("hello")
	}()
}

func execGoroutine_3() {
	sayHello := func() {
		fmt.Println("hello")
	}
	go sayHello()
}

func execGoroutine_4() {
	var wg sync.WaitGroup
	sayHello := func() {
		defer wg.Done()
		fmt.Println("hello")
	}
	wg.Add(1)
	go sayHello()
	wg.Wait()
}

func execGoroutine_5() {
	var wg sync.WaitGroup
	salutation := "helo"
	wg.Add(1)
	go func() {
		defer wg.Done()
		salutation = "welcome"
	}()
	wg.Wait()
	fmt.Println(salutation)
}

func execGoroutine_6() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
	}
	wg.Wait()
}

func execGoroutine_7() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait()
}

func execGoroutine_8() {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var ch <-chan interface{}
	var wg sync.WaitGroup
	noop := func() {
		wg.Done()
		// 計測のためにたくさんのゴルーチンをメモリに置いておきたいので、
		// 絶対に終了しないゴルーチンが必要です。
		// ここで何をしているか、いまはわからなくても問題ありません。
		// ただ、このゴルーチンはプロセスが終わるまで終了しないとだけ理解してください。
		<-ch // <1>
	}

	// ここで生成するゴルーチンの数を定義しています。
	// 対数の法則を使って漸近的にゴルーチンの数を増やしていきます。
	const numGoroutines = 1e4 // <2>
	wg.Add(numGoroutines)
	// ここでゴルーチン生成前のメモリ消費量を計測します。
	before := memConsumed() // <3>
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	// ここでゴルーチン生成後のメモリ消費量を計測します。
	after := memConsumed() // <4>
	fmt.Printf("%.3fkb\n", float64(after-before)/numGoroutines/1000)
}

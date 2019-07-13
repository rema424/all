package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"time"
)

func execTrace() {
	// トレース情報を出力するファイルを作成します。
	f, err := os.Create("trace.out") // <1>
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace output file: %v", err)
		}
	}()

	// トレースを開始し、トレース情報を<1>で作成したファイルに記録し始めます。
	// deferで関数を終了するときにトレースを停止していることに注意してください。
	if err := trace.Start(f); err != nil { // <2>
		panic(err)
	}
	defer trace.Stop()

	ctx := context.Background()
	// "makeCoffee"と名付けたTaskを作成しています。
	ctx, task := trace.NewTask(ctx, "makeCoffee") // <3>
	defer task.End()
	// "orderID"という名前を付けたLogに"1"というIDを渡しています。
	trace.Log(ctx, "orderID", "1") // <4>

	coffee := make(chan bool)

	go func() {
		trace.WithRegion(ctx, "extractCoffee", extractCoffee)
		coffee <- true
	}()
	<-coffee
}

func extractCoffee() {
	fmt.Println("extractCoffee running...")
	time.Sleep(5 * time.Second)
	fmt.Println("extractCoffee finished")
}

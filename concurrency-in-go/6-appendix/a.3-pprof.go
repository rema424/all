package main

import (
	"log"
	"os"
	"runtime/pprof"
	"time"
)

func execPprof() {
	log.SetFlags(log.Ltime | log.LUTC)
	log.SetOutput(os.Stdout)

	// 稼働しているゴルーチンの数を毎秒ログをとる
	go func() {
		goroutines := pprof.Lookup("goroutine")
		prof := newProfIfNotDef("goroutine")
		for range time.Tick(1 * time.Second) {
			log.Printf("goroutine count: %d\n", goroutines.Count())
			log.Printf("     prof count: %d\n", prof.Count())
		}
	}()

	// 決して終了しないゴルーチンをいくつか作成する
	var blockForever chan struct{}
	for i := 0; i < 10; i++ {
		go func() { <-blockForever }()
		time.Sleep(500 * time.Millisecond)
	}
}

func newProfIfNotDef(name string) *pprof.Profile {
	prof := pprof.Lookup(name)
	if prof == nil {
		prof = pprof.NewProfile(name)
	}
	return prof
}

package main

import (
	"context"
	"log"
	"os"
	"sync"

	"golang.org/x/time/rate"
)

func execSimpleRateLimit() {
	defer log.Printf("Done.")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := Open2()
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConnection.ReadFile(context.Background()); err != nil {
				log.Printf("cannot ReadFile: %v", err)
			}
			log.Printf("FeadFile")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConnection.ResolveAddress(context.Background()); err != nil {
				log.Printf("cannot ResolveAddress: %v", err)
			}
			log.Printf("ResolveAddress")
		}()
	}

	wg.Wait()
}

type APIConecction2 struct {
	rateLimiter *rate.Limiter
}

func Open2() *APIConecction2 {
	return &APIConecction2{
		// ここですべてのAPI接続に対して1秒に1つのイベントという流失制限をかけます。
		// 回復速度が1/1s（左）で、バケットの深さが1（右）
		rateLimiter: rate.NewLimiter(rate.Limit(1), 1), // <1>
	}
}

func (a *APIConecction2) ReadFile(ctx context.Context) error {
	// 流失制限のあとリクエストを完結させるのに十分な数のアクセストークンが揃うまで待機します。
	if err := a.rateLimiter.Wait(ctx); err != nil { // <2>
		return err
	}
	// 何かしたということにする
	return nil
}

func (a *APIConecction2) ResolveAddress(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil { // <2>
		return err
	}
	// 何かしたということにする
	return nil
}

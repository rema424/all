package main

import (
	"context"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func execMultiRateLimit() {
	defer log.Printf("Done.")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := Open3()
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot ReadFile: %v", err)
			}
			log.Printf("ReadFile")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot ResolveAddress: %v", err)
			}
			log.Printf("ResolveAddress")
		}()
	}

	wg.Wait()
}

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

// ここでRateLimiterインタフェースを定義してMultiLimiterが再帰的に
// 他のMultiLimiterインスタンスを定義できるようにします。
type RateLimiter interface { // <1>
	Wait(context.Context) error
	Limit() rate.Limit
}

type muliLimiter struct {
	limiters []RateLimiter
}

func MultiLimiter(limiters ...RateLimiter) *muliLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	// 最適化を実装して書くRateLimiterのLimit()でソートします。
	sort.Slice(limiters, byLimit) // <2>
	return &muliLimiter{limiters: limiters}
}

func (l *muliLimiter) Wait(ctx context.Context) error {
	// 1秒あたりの制限、1分あたりの制限、1日あたりの制限などすべての制限に対する流量を計算する。
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (l *muliLimiter) Limit() rate.Limit {
	// multiLimiterがインスタンス化されたときに子のRateLimiterインタフェースをソートするので、
	// 単純にスライスの最初の要素になっている、もっとも厳しい制限を返せます。
	return l.limiters[0].Limit() // <3>
}

type APIConnection3 struct {
	rateLimiter RateLimiter
}

func Open3() *APIConnection3 {
	// バースト性なしで1秒ごとの制限を定義します。
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1) // <1>
	// 1分ごとの制限をバースト性10に設定して、ユーザーに初期値のバッファを与えます。
	// 1秒ごとの制限によってシステムに対するリクエストで過負荷がかからないようにしています。
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10) // <2>
	return &APIConnection3{
		rateLimiter: MultiLimiter(secondLimit, minuteLimit),
	}
}

func (a *APIConnection3) ReadFile(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// 何かしたということにする
	return nil
}

func (a *APIConnection3) ResolveAddress(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}

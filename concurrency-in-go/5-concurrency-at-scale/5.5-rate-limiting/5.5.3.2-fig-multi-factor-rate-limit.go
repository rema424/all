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

func execMultiFactorRateLimit() {
	defer log.Printf("Done.")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := Open4()
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

func Per4(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

// ここでRateLimiter4インタフェースを定義してMultiLimiter4が再帰的に
// 他のMultiLimiter4インスタンスを定義できるようにします。
type RateLimiter4 interface { // <1>
	Wait(context.Context) error
	Limit() rate.Limit
}

type muliLimiter4 struct {
	limiters []RateLimiter4
}

func MultiLimiter4(limiters ...RateLimiter4) *muliLimiter4 {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	// 最適化を実装して書くRateLimiter4のLimit()でソートします。
	sort.Slice(limiters, byLimit) // <2>
	return &muliLimiter4{limiters: limiters}
}

func (l *muliLimiter4) Wait(ctx context.Context) error {
	// 1秒あたりの制限、1分あたりの制限、1日あたりの制限などすべての制限に対する流量を計算する。
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (l *muliLimiter4) Limit() rate.Limit {
	// multiLimiter4がインスタンス化されたときに子のRateLimiter4インタフェースをソートするので、
	// 単純にスライスの最初の要素になっている、もっとも厳しい制限を返せます。
	return l.limiters[0].Limit() // <3>
}

type APIConnection4 struct {
	networkLimit,
	diskLimit,
	apiLimit RateLimiter4
}

func Open4() *APIConnection4 {
	return &APIConnection4{
		// ここでAPI呼び出しの流量制限を設定しています。
		// 1秒間と1分間の両方のリクエスト数の上限を設定しています。
		apiLimit: MultiLimiter4( // <1>
			rate.NewLimiter(Per4(2, time.Second), 2),
			rate.NewLimiter(Per4(10, time.Minute), 10),
		),
		// ここでディスクの読み込みに対する流量制限を設けます。
		// 1秒間に1回の読み込みのみを制限とします。
		diskLimit: MultiLimiter4( // <2>
			rate.NewLimiter(rate.Limit(1), 1),
		),
		// ネットワークに関しては、1秒間に3リクエストまでの制限とします。
		networkLimit: MultiLimiter4( // <3>
			rate.NewLimiter(Per4(3, time.Second), 3),
		),
	}
}

func (a *APIConnection4) ReadFile(ctx context.Context) error {
	// ファイルの読み込みを行うときは、APIへの制限とディスクへの制限を組み合わせます。
	if err := MultiLimiter4(a.apiLimit, a.diskLimit).Wait(ctx); err != nil { // <4>
		return err
	}
	// 何かしたということにする
	return nil
}

func (a *APIConnection4) ResolveAddress(ctx context.Context) error {
	// ネットワークアクセスが必要な時は、APIへの制限とネットワークへの制限を組み合わせます。
	if err := MultiLimiter4(a.apiLimit, a.networkLimit).Wait(ctx); err != nil { // <5>
		return err
	}
	// Pretend we do work here
	return nil
}

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	// fmt.Println("Hello World")

	// {
	// 	var wg sync.WaitGroup
	// 	sayHello := func() {
	// 		defer wg.Done()
	// 		fmt.Println("hello")
	// 	}
	// 	wg.Add(1)
	// 	go sayHello()
	// 	wg.Wait()
	// }

	// {
	// 	var wg sync.WaitGroup
	// 	salutation := "hello"
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		salutation = "welcome"
	// 	}()
	// 	wg.Wait()
	// 	fmt.Println(salutation)
	// }

	// {
	// 	var wg sync.WaitGroup
	// 	for _, salutaion := range []string{"hello", "greetings", "good day"} {
	// 		wg.Add(1)
	// 		go func() {
	// 			defer wg.Done()
	// 			fmt.Println(salutaion)
	// 		}()
	// 	}
	// 	wg.Wait()
	// }

	// {
	// 	var wg sync.WaitGroup
	// 	for _, salutaion := range []string{"hello", "greetings", "good day"} {
	// 		wg.Add(1)
	// 		go func(salutaion string) {
	// 			defer wg.Done()
	// 			fmt.Println(salutaion)
	// 		}(salutaion)
	// 	}
	// 	wg.Wait()
	// }

	// {
	// 	memConsumed := func() uint64 {
	// 		runtime.GC()
	// 		var s runtime.MemStats
	// 		runtime.ReadMemStats(&s)
	// 		return s.Sys
	// 	}

	// 	var c <-chan interface{}
	// 	var wg sync.WaitGroup
	// 	noop := func() {
	// 		wg.Done()
	// 		<-c
	// 	}

	// 	const numGoroutines = 1e4
	// 	wg.Add(numGoroutines)
	// 	before := memConsumed()
	// 	for i := numGoroutines; i > 0; i-- {
	// 		go noop()
	// 	}
	// 	wg.Wait()
	// 	after := memConsumed()
	// 	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
	// }

	// context

	{
		locale := func(ctx context.Context) (string, error) {
			// select 文は処理をブロックする
			// context がキャンセルされるか、◯秒経つかするまでブロックする
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(3 * time.Second): // ここを 1, 2, 3 秒に書き換えて実行してみる
			}
			return "EN/US", nil
		}

		genGreeting := func(ctx context.Context) (string, error) {
			// ブロック処理のある locale() を呼び出している。
			// 親 context を元に新たな context を生成し、子関数に渡している。
			// 故に子関数のブロッキングを管理しているのは親ではなくこの関数である。

			ctx, cancel := context.WithTimeout(ctx, 2*time.Second) // 2 秒後に自動でキャンセルされる context
			// ブロッキングを管理する責任があるため、この関数の終わり際に context を確実に終了させる。
			// この context を用いて実行しているのは locale() だけで、親はキャンセルされない。
			defer cancel()

			// ここのコードは処理をブロックするわけではない
			// ブロックを発生させる処理の実装は locale() の中にある
			switch locale, err := locale(ctx); { // 要するにこの行で locale() から値が返却されるまで時間がかかる（中でブロッキングが発生している）
			case err != nil:
				return "", err
			case locale == "EN/US":
				return "hello", nil
			}
			return "", fmt.Errorf("unsupported locale")
		}

		printGreeting := func(ctx context.Context) error {
			// ブロック処理はない。
			// たとえ子の関数の中にブロッキングがあったとしても、
			// 親の context を渡しているため管理の責任は親にある。
			greeting, err := genGreeting(ctx)
			if err != nil {
				return err
			}
			fmt.Printf("%s world!\n", greeting)
			return nil
		}

		genFarewell := func(ctx context.Context) (string, error) {
			// ブロック処理のある locale() を呼び出している。
			// ただし、親から渡ってきた context をそのまま渡しているため、
			// この関数自体は子のブロッキングを管理する責任は持っていない。
			switch locale, err := locale(ctx); {
			case err != nil:
				return "", err
			case locale == "EN/US":
				return "goodbye", nil
			}
			return "", fmt.Errorf("unsupported locale")
		}

		printFarewell := func(ctx context.Context) error {
			// ブロック処理はない。
			// たとえ子の関数の中にブロッキングがあったとしても、
			// 親の context を渡しているため管理の責任は親にある。
			farewell, err := genFarewell(ctx)
			if err != nil {
				return err
			}
			fmt.Printf("%s world!\n", farewell)
			return nil
		}

		start := time.Now()

		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := printGreeting(ctx); err != nil {
				fmt.Printf("cannot print greeting: %v\n", err)
				// printGeering() でエラーが発生した場合は、
				// 同様の context を利用している全ての goroutine を解放する。
				// 他の goroutine が処理の途中だったとしても。
				// ただし、printGeering() 内でブロッキングが発生して、
				// error の返却がされない場合はここでずっと処理が止まることになる。cancel() もできない。
				cancel()
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := printFarewell(ctx); err != nil {
				fmt.Printf("cannot print farawell: %v\n", err)
				// printGeering() でエラーが発生した場合は、
				// 同様の context を利用している全ての goroutine を解放する。
				// 他の goroutine が処理の途中だったとしても。
				// ただし、printGeering() 内でブロッキングが発生して、
				// error の返却がされない場合はここでずっと処理が止まることになる。cancel() もできない。
				cancel()
			}
		}()

		wg.Wait()
		fmt.Println(time.Since(start))
	}
}

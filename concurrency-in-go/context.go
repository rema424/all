package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func execDoneChannel() {
	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)

	locale := func(done <-chan interface{}) (string, error) {
		select {
		case <-done:
			return "", fmt.Errorf("canceled")
		case <-time.After(1 * time.Second):
		}
		return "EN/US", nil
	}

	genGreeting := func(done <-chan interface{}) (string, error) {
		switch locale, err := locale(done); {
		case err != nil:
			return "", err
		case locale == "EN/US":
			return "hello", nil
		}
		return "", fmt.Errorf("unsupported locale")
	}

	printGreeting := func(done <-chan interface{}) error {
		greeting, err := genGreeting(done)
		if err != nil {
			return err
		}
		fmt.Printf("%s world!\n", greeting)
		return nil
	}

	genFarewell := func(done <-chan interface{}) (string, error) {
		switch locale, err := locale(done); {
		case err != nil:
			return "", err
		case locale == "EN/US":
			return "goodbye", nil
		}
		return "", fmt.Errorf("unsupported locale")
	}

	printFarewell := func(done <-chan interface{}) error {
		farewell, err := genFarewell(done)
		if err != nil {
			return err
		}
		fmt.Printf("%s world!\n", farewell)
		return nil
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Wait()
}

func execContext() {
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

func execDeadline() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	locale := func(ctx context.Context) (string, error) {
		if deadline, ok := ctx.Deadline(); ok {
			if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
				// genGreeting() で 30 秒のデッドラインを設定しているが、1 分未満のデッドラインはここで即時にエラーを発生させる。
				// genFarewell() はデッドラインの設定をしていないため、if をすり抜けて次の select でブロックされる。
				// genGreeting() はここでエラーとなる。
				// genFarewell() は select で1分のブロックをされてる間に genGreeting() の defer ctx.Done() が波及してすぐにエラー側に入る。
				return "", context.DeadlineExceeded
			}
		}

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(1 * time.Minute):
		}
		return "EN/US", nil
	}

	genGreeting := func(ctx context.Context) (string, error) {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second) // 1秒 -> デッドライン
		// ctx, cancel := context.WithTimeout(ctx, 6*time.Second) // 6秒 -> サクセス
		defer cancel()

		switch locale, err := locale(ctx); {
		case err != nil:
			return "", err
		case locale == "EN/US":
			return "hello", nil
		}
		return "", fmt.Errorf("unsupported locale")
	}

	printGreeting := func(ctx context.Context) error {
		greeting, err := genGreeting(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("%s world!\n", greeting)
		return nil
	}

	genFarewell := func(ctx context.Context) (string, error) {
		switch locale, err := locale(ctx); {
		case err != nil:
			return "", err
		case locale == "EN/US":
			return "good bye", nil
		}
		return "", fmt.Errorf("unsupported locale")
	}

	printFarewell := func(ctx context.Context) error {
		farewell, err := genFarewell(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("%s world!\n", farewell)
		return nil
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printGreeting(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v\n", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printFarewell(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v\n", err)
			cancel()
		}
	}()

	wg.Wait()
}

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func execValue(userID, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	handleResponse(ctx)
}

func handleResponse(ctx context.Context) {
	getUserID := func(c context.Context) string {
		return c.Value(ctxUserID).(string)
	}
	getAuthToken := func(c context.Context) string {
		return c.Value(ctxAuthToken).(string)
	}
	fmt.Printf(
		"handling response for %v (%v)",
		getUserID(ctx),
		getAuthToken(ctx),
	)
}

func execInterfaceKey() {
	type foo int
	type bar int
	m := make(map[interface{}]int)
	m[foo(1)] = 1
	m[bar(1)] = 2
	fmt.Printf("%v\n", m)
}

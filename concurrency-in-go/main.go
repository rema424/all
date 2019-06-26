package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// fmt.Println("Hello World")

	// 1.2.1 競合状態

	// {
	// 	var data int
	// 	go func() {
	// 		data++
	// 	}()
	// 	if data == 0 {
	// 		fmt.Printf("the value is %v.\n", data)
	// 	}
	// }

	// {
	// 	var data int
	// 	go func() {
	// 		data++
	// 	}()
	// 	time.Sleep(1 * time.Second)
	// 	if data == 0 {
	// 		fmt.Printf("the value is %v.\n", data)
	// 	}
	// }

	// 1.2.3 メモリアクセス同期

	// {
	// 	var data int
	// 	go func() {
	// 		data++
	// 	}()
	// 	if data == 0 {
	// 		fmt.Println("the value is 0.")
	// 	} else {
	// 		fmt.Printf("the value is %v.\n", data)
	// 	}
	// }

	// {
	// 	var memoryAccess sync.Mutex
	// 	var value int
	// 	go func() {
	// 		memoryAccess.Lock()
	// 		value++
	// 		memoryAccess.Unlock()
	// 	}()

	// 	memoryAccess.Lock()
	// 	if value == 0 {
	// 		fmt.Printf("the value is %v.\n", value)
	// 	} else {
	// 		fmt.Printf("the value is %v.\n", value)
	// 	}
	// 	memoryAccess.Unlock()
	// }

	// 1.2.4 デッドロック、ライブロック、リソース枯渇

	// デッドロック
	// {
	// 	type value struct {
	// 		mu    sync.Mutex
	// 		value int
	// 	}

	// 	var wg sync.WaitGroup
	// 	printSum := func(v1, v2 *value) {
	// 		defer wg.Done()
	// 		v1.mu.Lock()
	// 		defer v1.mu.Unlock()

	// 		time.Sleep(2 * time.Second)
	// 		v2.mu.Lock()
	// 		defer v2.mu.Unlock()

	// 		fmt.Printf("sum=%v\n", v1.value+v2.value)
	// 	}

	// 	var a, b value
	// 	wg.Add(2)
	// 	go printSum(&a, &b)
	// 	go printSum(&b, &a) // fatal error: all goroutines are asleep - deadlock!
	// 	wg.Wait()
	// }

	// ライブロック
	// {
	// 	cadence := sync.NewCond(&sync.Mutex{})
	// 	go func() {
	// 		for range time.Tick(1 * time.Millisecond) {
	// 			cadence.Broadcast()
	// 		}
	// 	}()

	// 	takeStep := func() {
	// 		cadence.L.Lock()
	// 		cadence.Wait()
	// 		cadence.L.Unlock()
	// 	}

	// 	tryDir := func(dirName string, dir *int32, out *bytes.Buffer) bool {
	// 		fmt.Fprintf(out, " %v", dirName)
	// 		atomic.AddInt32(dir, 1)
	// 		takeStep()
	// 		if atomic.LoadInt32(dir) == 1 {
	// 			fmt.Fprint(out, ". Success!")
	// 			return true
	// 		}
	// 		takeStep()
	// 		atomic.AddInt32(dir, -1)
	// 		return false
	// 	}

	// 	var left, right int32
	// 	tryLeft := func(out *bytes.Buffer) bool {
	// 		return tryDir("left", &left, out)
	// 	}
	// 	tryRight := func(out *bytes.Buffer) bool {
	// 		return tryDir("right", &right, out)
	// 	}
	// 	walk := func(walking *sync.WaitGroup, name string) {
	// 		var out bytes.Buffer
	// 		defer func() {
	// 			fmt.Println(out.String())
	// 		}()
	// 		defer walking.Done()
	// 		fmt.Fprintf(&out, "%v is trying to scoot: ", name)
	// 		for i := 0; i < 5; i++ {
	// 			if tryLeft(&out) || tryRight(&out) {
	// 				return
	// 			}
	// 		}
	// 		fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name)
	// 	}

	// 	var peopleInHallway sync.WaitGroup
	// 	peopleInHallway.Add(2)
	// 	go walk(&peopleInHallway, "Alice")
	// 	go walk(&peopleInHallway, "Barbara")
	// 	peopleInHallway.Wait()
	// }

	// リソース枯渇

	// ----------------------
	// 4.12 context パッケージ
	// ----------------------

	// done チャネル

	{
		var wg sync.WaitGroup
		done := make(chan interface{})
		defer close(done)

		locale := func(done <-chan interface{}) (string, error) {
			select {
			case <-done:
				return "", fmt.Errorf("canceled")
			case <-time.After(1 * time.Minute):
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

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := printGreeting(done); err != nil {
				fmt.Printf("%v", err)
				return
			}
		}()
	}

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
}

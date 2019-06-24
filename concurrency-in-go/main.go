package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	fmt.Println("Hello World")

	{
		var wg sync.WaitGroup
		sayHello := func() {
			defer wg.Done()
			fmt.Println("hello")
		}
		wg.Add(1)
		go sayHello()
		wg.Wait()
	}

	{
		var wg sync.WaitGroup
		salutation := "hello"
		wg.Add(1)
		go func() {
			defer wg.Done()
			salutation = "welcome"
		}()
		wg.Wait()
		fmt.Println(salutation)
	}

	{
		var wg sync.WaitGroup
		for _, salutaion := range []string{"hello", "greetings", "good day"} {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Println(salutaion)
			}()
		}
		wg.Wait()
	}

	{
		var wg sync.WaitGroup
		for _, salutaion := range []string{"hello", "greetings", "good day"} {
			wg.Add(1)
			go func(salutaion string) {
				defer wg.Done()
				fmt.Println(salutaion)
			}(salutaion)
		}
		wg.Wait()
	}

	{
		memConsumed := func() uint64 {
			runtime.GC()
			var s runtime.MemStats
			runtime.ReadMemStats(&s)
			return s.Sys
		}

		var c <-chan interface{}
		var wg sync.WaitGroup
		noop := func() {
			wg.Done()
			<-c
		}

		const numGoroutines = 1e4
		wg.Add(numGoroutines)
		before := memConsumed()
		for i := numGoroutines; i > 0; i-- {
			go noop()
		}
		wg.Wait()
		after := memConsumed()
		fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
	}
}

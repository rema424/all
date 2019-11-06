package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx1 := context.Background()
	ctx1_2 := context.WithValue(ctx1, "foo", "bar")
	ctx1_2_2 := context.WithValue(ctx1_2, "hoge", "fuga")
	ctx1_2_3 := context.WithValue(ctx1_2, "foo", "buzz")
	fmt.Println("========================================")
	fmt.Printf("%s\t%v\n", "ctx1", ctx1)
	fmt.Printf("%s\t%v\n", "ctx1_2", ctx1_2)
	fmt.Printf("%s\t%v\n", "ctx1_2_2", ctx1_2_2)
	fmt.Printf("%s\t%v\n", "ctx1_2_3", ctx1_2_3)
	fmt.Println("========================================")
	fmt.Printf("%s\t%v\n", "ctx1:foo", ctx1.Value("foo"))
	fmt.Printf("%s\t%v\n", "ctx1_2:foo", ctx1_2.Value("foo"))
	fmt.Printf("%s\t%v\n", "ctx1_2_2:foo", ctx1_2_2.Value("foo"))
	fmt.Printf("%s\t%v\n", "ctx1_2_3:foo", ctx1_2_3.Value("foo"))
	fmt.Println("========================================")
	fmt.Printf("%s\t%v\n", "ctx1:hoge", ctx1.Value("hoge"))
	fmt.Printf("%s\t%v\n", "ctx1_2:hoge", ctx1_2.Value("hoge"))
	fmt.Printf("%s\t%v\n", "ctx1_2_2:hoge", ctx1_2_2.Value("hoge"))
	fmt.Printf("%s\t%v\n", "ctx1_2_3:hoge", ctx1_2_3.Value("hoge"))
	fmt.Println("========================================")

	// ctx1_3, cancel1_3 := context.WithTimeout(ctx1, 5*time.Second)
	// _, cancel1_3 := context.WithTimeout(ctx1, 5*time.Second)
	ctx1_3, _ := context.WithTimeout(ctx1, 5*time.Second)
	// time.AfterFunc(3*time.Second, cancel1_3)
	fmt.Println(ctx1_3)

	go func() {
		select {
		case <-ctx1_3.Done():
			fmt.Println("ctx1_3", ctx1_3.Err())
		}
	}()

	for {
		select {
		case <-ctx1_2.Done():
			fmt.Println("ctx1_2", ctx1_2.Err())
		default:
			fmt.Println("waiting...")
			time.Sleep(time.Second)
		}
	}
}

// func main() {
// 	fmt.Println("========================================")
// 	ctx := context.Background()
// 	ctx2 := context.WithValue(ctx, "foo", "bar")
// 	fmt.Println("ctx", ctx)
// 	fmt.Println("ctx2", ctx2)
// 	fmt.Println("========================================")
// 	fmt.Println("ctx:foo", ctx.Value("foo"))
// 	fmt.Println("ctx2:foo", ctx2.Value("foo"))
// 	fmt.Println("========================================")
// 	ctx3 := context.WithValue(ctx, "hoge", "fuga")
// 	fmt.Println("ctx", ctx)
// 	fmt.Println("ctx2", ctx2)
// 	fmt.Println("ctx3", ctx3)
// 	fmt.Println("========================================")
// 	fmt.Println("ctx:foo", ctx.Value("foo"))
// 	fmt.Println("ctx2:foo", ctx2.Value("foo"))
// 	fmt.Println("ctx3:foo", ctx3.Value("foo"))
// 	fmt.Println("========================================")
// 	fmt.Println("ctx:hoge", ctx.Value("hoge"))
// 	fmt.Println("ctx2:hoge", ctx2.Value("hoge"))
// 	fmt.Println("ctx3:hoge", ctx3.Value("hoge"))
// 	fmt.Println("========================================")
// 	ctx2_2, cancel2_2 := context.WithTimeout(ctx2, 10*time.Second)
// 	ctx2_3, cancel2_3 := context.WithDeadline(ctx2_2, time.Now().Add(15*time.Second))
// 	ctx2_4 := context.WithValue(ctx2_3, "aaa", "bbb")
// 	defer cancel2_2()
// 	defer cancel2_3()
// 	fmt.Println("ctx2", ctx2)
// 	fmt.Println("ctx2_2", ctx2_2)
// 	fmt.Println("ctx2_3", ctx2_3)
// 	fmt.Println("ctx2_4", ctx2_4)
// 	fmt.Println("ctx2_3:foo", ctx2_3.Value("foo"))
// 	fmt.Println("ctx2_3:aaa", ctx2_3.Value("aaa"))
// 	fmt.Println("ctx2_4:foo", ctx2_4.Value("foo"))
// 	fmt.Println("ctx2_4:aaa", ctx2_4.Value("aaa"))

// 	loop := func(ctx context.Context, id int) {
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				fmt.Println(id, ctx.Err())
// 				return
// 			default:
// 				fmt.Println(id, rand.Int())
// 				time.Sleep(time.Second)
// 			}
// 		}
// 	}

// 	// ctx2WithTimeout, cancel2 := context.WithTimeout(ctx2, 5*time.Second)
// 	// fmt.Println("ctx2:foo", ctx2.Value("foo"))
// 	// fmt.Println("ctx2WithTimeout:foo", ctx2WithTimeout.Value("foo"))

// 	// defer cancel2()

// 	time.AfterFunc(9*time.Second, cancel2_2)
// 	go loop(ctx2_2, 1)
// 	// go loop(ctx2WithTimeout, 2)
// 	// go loop(ctx3, 3)

// 	select {
// 	case <-ctx2_3.Done():
// 		fmt.Println(ctx2_3.Err())
// 	}
// }

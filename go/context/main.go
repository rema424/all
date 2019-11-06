package main

import (
	"context"
	"fmt"
)

func main() {
	fmt.Println("========================================")
	ctx := context.Background()
	ctx2 := context.WithValue(ctx, "foo", "bar")
	fmt.Println("ctx", ctx)
	fmt.Println("ctx2", ctx2)
	fmt.Println("========================================")
	fmt.Println("ctx:foo", ctx.Value("foo"))
	fmt.Println("ctx2:foo", ctx2.Value("foo"))
	fmt.Println("========================================")
	ctx3 := context.WithValue(ctx, "hoge", "fuga")
	fmt.Println("ctx", ctx)
	fmt.Println("ctx2", ctx2)
	fmt.Println("ctx3", ctx3)
	fmt.Println("========================================")
	fmt.Println("ctx:foo", ctx.Value("foo"))
	fmt.Println("ctx2:foo", ctx2.Value("foo"))
	fmt.Println("ctx3:foo", ctx3.Value("foo"))
	fmt.Println("========================================")
	fmt.Println("ctx:hoge", ctx.Value("hoge"))
	fmt.Println("ctx2:hoge", ctx2.Value("hoge"))
	fmt.Println("ctx3:hoge", ctx3.Value("hoge"))
	fmt.Println("========================================")
}

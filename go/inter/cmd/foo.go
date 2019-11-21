package main

import (
	"inter/internal/hello"
	qux "inter/qux/internal/hello"
)

func main() {
	hello.Hello()
	qux.Hello()
}

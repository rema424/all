package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	typeName := flag.String("type", "", "todo usage")
	flag.Parse()
	if len(*typeName) == 0 {
		flag.Usage()
		os.Exit(2)
	} else {
		fmt.Printf("%-10s\n", *typeName)
	}
}

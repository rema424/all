package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"

	"github.com/k0kubun/pp"
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
	fmt.Println(
		build.AllowBinary,
		build.FindOnly,
		build.IgnoreVendor,
		build.ImportComment,
	)
	// p, err := build.Default.ImportDir("greet", build.AllowBinary)
	// p, err := build.Default.ImportDir("greet", build.FindOnly)
	// p, err := build.Default.ImportDir("greet", build.IgnoreVendor)
	p, err := build.Default.ImportDir("greet", build.ImportComment)
	if err != nil {
		panic(err)
	}
	pp.Println(p)
}

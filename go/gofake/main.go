package main

import (
	"flag"
	"gofake/gofake"
	"os"
)

var (
	typ string
	dir string
)

func init() {
	flag.StringVar(&typ, "type", "Greeter", "type name")
	flag.StringVar(&dir, "dir", ".", "dir name")
	flag.Parse()

	if len(typ) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	// b, err := gofake.MakeFileContents("animal")
	// fmt.Println(string(b), err)
	gofake.Run(typ, dir)
}

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
	flag.StringVar(&typ, "type", ".", "type name")
	flag.StringVar(&dir, "dir", ".", "dir name")
	flag.Parse()

	if len(typ) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	gofake.Run(typ, dir)
}

package gofake

import (
	"fmt"
	"go/build"

	"github.com/k0kubun/pp"
)

func Run(typ, dir string) {
	fmt.Println(
		build.AllowBinary,
		build.FindOnly,
		build.IgnoreVendor,
		build.ImportComment,
	)
	// p, err := build.Default.ImportDir("greet", build.AllowBinary)
	// p, err := build.Default.ImportDir("greet", build.FindOnly)
	// p, err := build.Default.ImportDir("greet", build.IgnoreVendor)
	pkg, err := build.Default.ImportDir(dir, build.ImportComment)
	if err != nil {
		panic(err)
	}
	files := make([]string, len(pkg.GoFiles))
	for i, f := range pkg.GoFiles {
		files[i] = f
	}
	pp.Println(pkg)
	fmt.Printf("%q\n", files)
}

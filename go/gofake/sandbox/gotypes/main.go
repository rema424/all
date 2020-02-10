package main

import (
	"fmt"
	"log"
	"path/filepath"

	"golang.org/x/tools/imports"
)

func main() {
	srcPath := filepath.Join(".", "__tmp__.go")
	structPath := "hello.Hello"

	src := []byte("package hack\n" + "var i " + structPath)
	imp, err := imports.Process(srcPath, src, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(imp))

	// conf := types.Config{Importer: importer.Default()}

	// pkg, err := conf.Check("hello", fset, []*ast.File{f}, nil)
	// if err != nil {
	// 	log.Fatal(err) // type error
	// }
	// obj := pkg.Scope().Lookup("Hello")

	// fmt.Println(obj)
	// fmt.Println("exported:", obj.Exported())
	// fmt.Println("id:", obj.Id())
	// fmt.Println("name:", obj.Name())
	// fmt.Println("parent:", obj.Parent())
	// fmt.Println("pkg:", obj.Pkg())
	// fmt.Println("pos:", obj.Pos())
	// fmt.Println("string:", obj.String())
	// fmt.Println("type:", obj.Type())
	// pp.Println(obj)

	// fmt.Printf("Package  %q\n", pkg.Path())
	// fmt.Printf("Name:    %s\n", pkg.Name())
	// fmt.Printf("Imports: %s\n", pkg.Imports())
	// fmt.Printf("Scope:   %s\n", pkg.Scope())
}

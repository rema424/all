package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"log"

	"golang.org/x/tools/go/loader"
)

var source = `
package person

type Person struct {
	ID      int64
	Profile Profile
	Address *Address
	Hobbies []Hobby
	Pets    []*Pet
}

type Profile struct {
	Name string
}

type Address struct {
	Zip string
}

type Hobby struct {
	Name string
}

type Pet struct {
	Name string
}
`

func main() {
	fmt.Println(source)

	ldr := loader.Config{ParserMode: parser.ParseComments}
	f, err := ldr.ParseFile("tmp.go", source)
	if err != nil {
		log.Fatal(err)
	}

	ldr.CreateFromFiles("person", f)

	program, err := ldr.Load()
	if err != nil {
		log.Fatal(err)
	}
	// pp.Println(program)
	// for _, v := range program.AllPackages {
	// 	pp.Println(v)
	// }
	person := program.Package("person")
	// person.Info.
	for _, f := range person.Files {
		ast.Inspect(f, func(node ast.Node) bool {
			return true
		})
	}
}

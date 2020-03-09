package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"strings"
)

const src = `
package person

type Person struct {
	ID      int64
	Profile Profile
	Address *Address
	Hobbies []Hobby
	Pets    []*Pet
	Job     *struct {
		Name string
	}
	Food struct {
		Name string
	}
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
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "tmp.go", src, parser.Mode(0))

	ast.Inspect(f, func(n ast.Node) bool {
		// *ast.TypeSpecに当たるまで再帰的に走査
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		// *ast.StructTypeでない場合は走査打ち切り
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return false
		}
		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				fmt.Print(name.Name, " ")
				var b strings.Builder
				if err := printer.Fprint(&b, fset, field.Type); err != nil {
					log.Fatal(err)
				}
				fmt.Print(b.String())
				// if star, ok := field.Type.(*ast.StarExpr); !ok {
				// 	fmt.Print(field.Type)
				// } else {
				// 	fmt.Print("*")
				// 	if sel, ok := star.X.(*ast.SelectorExpr); !ok {
				// 		// pp.Println(star)
				// 		fmt.Print(field.Type)
				// 	} else {
				// 		fmt.Print(sel.X)
				// 		fmt.Print(".")
				// 		fmt.Print(sel.Sel)
				// 	}
				// }
				fmt.Print("\n")
			}
		}
		// printer.Fprint(os.Stdout, fset, typeSpec)
		// pp.Println(structType)
		return false
	})
}

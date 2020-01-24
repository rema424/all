package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "sample/person.go", nil, parser.Mode(0))

	// for _, d := range f.Decls {
	// 	ast.Print(fset, d)
	// 	fmt.Println()
	// }

	// fmt.Println(f.Doc)
	// fmt.Println(f.Package)
	// fmt.Println(f.Name)
	// fmt.Println(f.Decls)
	// fmt.Println(f.Imports)
	// fmt.Println(f.Comments)
	// ast.Print(fset, f)

	// for _, decl := range f.Decls {
	// 	switch gen := decl.(type) {
	// 	case *ast.GenDecl:
	// 		for _, spec := range gen.Specs {
	// 			if typspec, ok := spec.(*ast.TypeSpec); ok {
	// 				if strct, ok := typspec.Type.(*ast.StructType); ok {
	// 					fmt.Println(typspec.Name)
	// 					for _, f := range strct.Fields.List {
	// 						for _, name := range f.Names {
	// 							fmt.Print(name.String(), " ")
	// 						}
	// 						fmt.Println(f.Type, f.Tag)
	// 					}
	// 					fmt.Println(strct.Fields)
	// 					ast.Print(fset, decl)
	// 					fmt.Println()
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	ast.Inspect(f, func(node ast.Node) bool {
		spec, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}
		fmt.Println(spec.Name)

		ast.Inspect(spec, func(node ast.Node) bool {
			switch node := node.(type) {
			case (*ast.StructType):
				for _, f := range node.Fields.List {
					for _, name := range f.Names {
						fmt.Print(name.String(), " ")
					}
					fmt.Println(f.Type)
					ast.Print(fset, f)
				}
			case (*ast.InterfaceType):
				for _, f := range node.Methods.List {
					for _, name := range f.Names {
						fmt.Print(name.String(), " ")
					}
					fmt.Println(f.Type)
					ast.Print(fset, f)
				}
			}
			return true
		})

		return true
	})
}

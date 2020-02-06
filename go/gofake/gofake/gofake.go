package gofake

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func Run(typ, dir string) {
	// fmt.Println(
	// 	build.AllowBinary,
	// 	build.FindOnly,
	// 	build.IgnoreVendor,
	// 	build.ImportComment,
	// )
	// p, err := build.Default.ImportDir("greet", build.AllowBinary)
	// p, err := build.Default.ImportDir("greet", build.FindOnly)
	// p, err := build.Default.ImportDir("greet", build.IgnoreVendor)
	pkg, err := build.Default.ImportDir(dir, build.ImportComment)
	if err != nil {
		panic(err)
	}
	paths := make([]string, len(pkg.GoFiles))
	for i, f := range pkg.GoFiles {
		paths[i] = filepath.Join(pkg.Dir, f)
	}
	// pp.Println(pkg)
	fmt.Printf("%q\n", paths)

	fset := token.NewFileSet()
	for _, p := range paths {
		f, err := parser.ParseFile(fset, p, nil, 0)
		if err != nil {
			panic(err)
		}

		ast.Inspect(f, func(node ast.Node) bool {
			typeSpec, ok := node.(*ast.TypeSpec)
			if !ok {
				return true
			}
			if typeSpec.Name.String() != typ {
				return true
			}

			ast.Inspect(typeSpec, func(node ast.Node) bool {
				structType, ok := node.(*ast.StructType)
				if !ok {
					return true
				}

				// ----------------------------
				// Struct Pcoccessing Start
				// ----------------------------
				fmt.Println(p)
				fmt.Println(filepath.Base(p))
				fmt.Println(typeSpec.Name)

				fmt.Println()
				genFilePath := makeGenFilePath(p)
				isGenFileExists := isFileExists(genFilePath)
				if isGenFileExists {
					modifyGenFile()
				} else {
					makeGenFile(genFilePath)
				}

				for _, field := range structType.Fields.List {
					fmt.Print(field.Type, " ")
					for i, name := range field.Names {
						if i == len(field.Names)-1 {
							fmt.Println(name)
						} else {
							fmt.Print(name, " ")
						}
					}
				}

				// ----------------------------
				// Struct Pcoccessing End
				// ----------------------------

				return false
			})
			return false
		})
		// pp.Println(f)
	}
}

func modifyGenFile() {
	fmt.Println("modifyGenFile")
}

func makeGenFile(filePath string) error {
	_, err := os.Create(filePath)
	if err != nil {
		return err
	}
	return nil
}

func MakeFileContents(pkg string) ([]byte, error) {
	var b bytes.Buffer
	b.WriteString("package ")
	b.WriteString(pkg)
	b.WriteString("\n")

	return format.Source(b.Bytes())
}

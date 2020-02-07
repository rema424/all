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

		structs := ExtractStructs(f)
		fmt.Println(structs)
		b := MakeFileContents(pkg.Name, structs)
		src, err := format.Source(b)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
		fmt.Println(string(src))
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

func MakeFileContents(pkg string, ss []Struct) []byte {
	var b bytes.Buffer
	b.WriteString("package ")
	b.WriteString(pkg)
	b.WriteString("\n")

	for _, s := range ss {
		b.WriteString(fmt.Sprintf("type %ss []*%s\n", s.Name, s.Name))
		for _, f := range s.Fields {
			b.WriteString(fmt.Sprintf("func (ss %ss) %ss() []%s { ", s.Name, f.Name, f.Type))
			b.WriteString(fmt.Sprintf("res := make([]%s, len(ss)); ", f.Type))
			b.WriteString("for i, s := range ss { ")
			b.WriteString("res[i] = s." + f.Name + " ")
			b.WriteString("}; ")
			b.WriteString("return res ")
			b.WriteString("}\n\n")
		}
	}

	return b.Bytes()
}

func ExtractStructs(f *ast.File) []Struct {
	ss := make([]Struct, 0, 10)
	ast.Inspect(f, func(node ast.Node) bool {
		typeSpec, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}

		ast.Inspect(typeSpec, func(node ast.Node) bool {
			structType, ok := node.(*ast.StructType)
			if !ok {
				return true
			}

			fields := make([]Field, 0, 10)
			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					fields = append(fields, Field{
						Name: name.Name,
						Type: fmt.Sprint(field.Type),
					})
				}

				ss = append(ss, Struct{
					Name:   typeSpec.Name.Name,
					Fields: fields,
				})
			}
			return false
		})
		return false
	})
	return ss
}

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
}

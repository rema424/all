package gofake

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	// path, err := filepath.Abs(dir)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(path)
	pkg, err := build.Default.ImportDir(dir, build.ImportComment)
	// pkg, err := build.Default.ImportDir("greeter/", build.ImportComment)
	if err != nil {
		panic(err)
	}
	paths := make([]string, 0, len(pkg.GoFiles))
	for _, f := range pkg.GoFiles {
		if strings.Contains(f, "_gen") {
			continue
		}
		path := filepath.Join(pkg.Dir, f)
		paths = append(paths, path)
	}
	// pp.Println(pkg)
	// fmt.Printf("%q\n", paths)

	fset := token.NewFileSet()
	for _, p := range paths {
		f, err := parser.ParseFile(fset, p, nil, 0)
		if err != nil {
			panic(err)
		}

		structs := ExtractStructs(fset, f, typ)
		// pp.Println(structs)
		b := MakeFileContents(pkg.Name, structs)
		src, err := format.Source(b)
		if err != nil {
			panic(err)
		}
		// fmt.Println(string(b))
		// fmt.Println(string(src))
		file, err := os.OpenFile(makeGenFilePath(p), os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		if _, err := file.Write(src); err != nil {
			panic(err)
		}
		file.Close()
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

func ExtractStructs(fset *token.FileSet, f *ast.File, typ string) []Struct {
	ss := make([]Struct, 0, 10)

	ast.Inspect(f, func(node ast.Node) bool {
		// *ast.TypeSpecに当たるまで再帰的に走査
		typeSpec, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}
		// 対象でない場合は走査打ち切り
		if typeSpec.Name.Name != typ {
			return false
		}
		// *ast.StructTypeでない場合は走査打ち切り
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return false
		}

		fields := make([]Field, 0, 10)
		// fmt.Println("go for loop")
		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				var b strings.Builder
				if err := printer.Fprint(&b, fset, field.Type); err != nil {
					log.Fatal(err)
				}
				fields = append(fields, Field{
					Name: name.Name,
					Type: b.String(),
				})
			}
		}
		s := Struct{
			Name:   typeSpec.Name.Name,
			Fields: fields,
		}

		ss = append(ss, s)
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

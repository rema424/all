package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"

	"github.com/k0kubun/pp"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/buildutil"
	"golang.org/x/tools/go/packages"
	"golang.org/x/xerrors"
)

var errNotFound = xerrors.New("no struct literal found at selection")

var (
	filename string
	modified bool
	offset   int
	line     int
)

func init() {
	flag.StringVar(&filename, "file", "", "filename")
	flag.BoolVar(&modified, "modified", false, "read an archive of modified files from stdin")
	flag.IntVar(&offset, "offset", 0, "byte offset of the struct literal, optonal if -offset is present")
	flag.IntVar(&line, "line", 0, "line number of the struct literal, optional if -line is present")
	flag.Parse()

	if (offset == 0 && line == 0) || filename == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("fillstruct: ")
	fmt.Println(filename, modified, offset, line)

	path, err := absPath(filename)
	if err != nil {
		log.Fatal(err)
	}

	var overlay map[string][]byte
	if modified {
		overlay, err = buildutil.ParseOverlayArchive(os.Stdin)
		if err != nil {
			log.Fatalf("invalid archive: %v", err)
		}
	}

	pp.Println("===================")
	cfg := &packages.Config{
		Overlay: overlay,
		Mode:    packages.LoadAllSyntax,
		Tests:   true,
		Dir:     filepath.Dir(path),
		Fset:    token.NewFileSet(),
		Env:     os.Environ(),
	}

	pkgs, err := packages.Load(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(path)
	byOffset(pkgs, path, offset)
	// pp.Println("pkgs", pkgs)
}

func absPath(filename string) (string, error) {
	eval, err := filepath.EvalSymlinks(filename)
	if err != nil {
		return "", err
	}
	return filepath.Abs(eval)
}

func byOffset(lprog []*packages.Package, path string, offset int) error {
	f, pkg, pos, err := findPos(lprog, path, offset)
	if err != nil {
		return err
	}
	fmt.Println(f, pkg, pos)
	return nil
}

func findPos(lprog []*packages.Package, path string, off int) (*ast.File, *packages.Package, token.Pos, error) {
	for _, pkg := range lprog {
		pp.Println("pkg", pkg)
		astFiles := pkg.Syntax
		for _, astFile := range astFiles {
			if tokenFile := pkg.Fset.File(astFile.Pos()); tokenFile.Name() == path {
				pp.Println("astFile", astFile)
				pp.Println("tokenFile", tokenFile)
				if off > tokenFile.Size() {
					return nil, nil, 0, fmt.Errorf("file size (%d) is smaller than given offset (%d)", tokenFile.Size(), off)
				}
				return astFile, pkg, tokenFile.Pos(off), nil
			}
		}
	}
	return nil, nil, 0, fmt.Errorf("could not find file %q", path)
}

// litInfo contains the information about
// a literal to fill with zero values.
type litInfo struct {
	typ       types.Type   // the base type of the literal
	name      *types.Named // name of the type or nil, e.g. for an anonymous struct type
	hideType  bool         // flag to hide the element type inside an array, slice or map literal
	isPointer bool         // true if the literal is of a pointer type
}

func findCompositeLit(f *ast.File, info *types.Info, pos token.Pos) (*ast.CompositeLit, litInfo, error) {
	var linfo litInfo
	path, _ := astutil.PathEnclosingInterval(f, pos, pos)
	pp.Println(path)
	return nil, linfo, errNotFound
}

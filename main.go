package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	flagFunc			= flag.String("f", "", "Function to analyse")
	flagFuncPrefix 		= flag.String("p", "", "Prefix to keep function")
	flagArgumentType 	= flag.String("arg", "", "Type of the argument of target functions")
)

func OtherConditions() {
	return
}

func GetPackage() *packages.Package {
	if flag.NArg() > 1 {
		panic("Should not receive more than one argument")
	}
	pkg_name := "."
	if flag.NArg() == 1 {
		pkg_name = flag.Arg(0)
	}
	fmt.Printf("Loading package %q\n", pkg_name)

	cfg := &packages.Config{Mode: packages.LoadAllSyntax, Tests: true}
	pkgs, err := packages.Load(cfg, pkg_name)
	if err != nil {
		panic(fmt.Sprintf("Package %v not found", flag.Args()))
	}
	if len(pkgs) != 1 {
		paths := make([]string, len(pkgs))
		for _, p := range pkgs {
			paths = append(paths, p.PkgPath)
		}
		panic(fmt.Sprintf("Cannot build multiple packages, %q resolved to %v", pkg_name, strings.Join(paths, ", ")))
	}

	pkg := pkgs[0]
	for _, name := range pkg.CompiledGoFiles {
		fmt.Printf("\t- %s\n", name)
	}
	return pkg
}

type FunctionExtractor struct {
	f *ast.File
	p *packages.Package
	tgt_func *[]*ast.FuncDecl
}

func IsFuncQualified(f *ast.FuncDecl, p *packages.Package) bool {
	info := p.TypesInfo
	f_name := f.Name.Name

	obj, ok := info.Defs[f.Name]
	if !ok {
		panic(fmt.Sprintf("Object of %s is not found", f.Name.Name))
	}

	obj_f, ok := obj.(*types.Func)
	if !ok {
		panic(fmt.Sprintf("%s cannot be converted into *types.Func", obj.String()))
	}

	sig, ok := obj_f.Type().Underlying().(*types.Signature)
	if !ok {
		panic(fmt.Sprintf("%s cannot be converted into *types.Signature", obj_f.String()))
	}

	// Filter out methods
	if sig.Recv() != nil {
		return false
	}
	// Filter out variadic functions, e.g. func Printf(a int, b ...interface{})
	if sig.Variadic() {
		return false
	}
	// Filter out functions with undisirable prefix.
	// As expectation, we only keep functions with "Test" as prefix
	if !strings.HasPrefix(f_name, *flagFuncPrefix) {
		return false
	}

	if !(sig.Params().Len() == 1 && sig.Results().Len() == 0) {
		return false
	}
	n := sig.Params().At(0).Type().String()
	if !strings.Contains(n, *flagArgumentType) {
		return false
	}

	// for i := 0; i < sig.Params().Len(); i++{
	// 	switch n := sig.Params().At(i).Type().(type) {
	// 	case *types.Named:
	// 		fmt.Printf("\tNamed types, %s\n", n.String())
	// 	case *types.Basic:
	// 		fmt.Printf("\tBasic types\n")
	// 	case *types.Interface:
	// 		fmt.Printf("\tInterface types\n")
	// 	case *types.Pointer:
	// 		fmt.Printf("\tPointer types, ")
	// 		elem := n.Elem()
	// 		fmt.Printf("%s", elem)
	// 		switch elem_t := elem.(type) {
	// 		case *types.Named:
	// 			tn := elem_t.Obj()
	// 			fmt.Printf("%s(%s)", tn.Name(), tn.Id())
	// 		case *types.Basic:
	// 			OtherConditions()
	// 		}
	// 		fmt.Printf("\n")
	// 	default:
	// 		fmt.Printf("\tOther types\n")
	// 		return false
	// 	}
	// }

	fmt.Printf("- %q found, %s\n", obj_f.Name(), sig.String())
	return true
}

func (fe *FunctionExtractor) Visit (node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if IsFuncQualified(n, fe.p) {
			*fe.tgt_func = append(*fe.tgt_func, n)
		}
		return nil
	}
	return fe
}

func ExtractFunctionFromFile(f *ast.File,
							p *packages.Package,
							tgt_func *[]*ast.FuncDecl) {
	fe := &FunctionExtractor{f: f, p: p, tgt_func: tgt_func}
	ast.Walk(fe, f)
}

func AnalyseFunction(f_ast *ast.FuncDecl, pkg *packages.Package) {
	oracle := Oracle {f_ast: f_ast,
					 fset: pkg.Fset,
					 pkg: pkg}
	_ = oracle

	ast.Walk(&oracle, f_ast)
}

func ExtractFunctionFromPackage() (result []*ast.FuncDecl) {
	tgt_pkg := GetPackage()
	var tgt_func []*ast.FuncDecl

	for i, fullname := range tgt_pkg.CompiledGoFiles {
		fname := filepath.Base(fullname)
		if !strings.HasSuffix(fname, ".go") {
			// This is a cgo-generated file.
			// Currently does not work.
			// We copied the original Go file as part of copyPackageRewrite,
			// so we can just skip this one.
			// See https://golang.org/issue/30479.
			continue
		}
		f := tgt_pkg.Syntax[i]
		ExtractFunctionFromFile(f, tgt_pkg, &tgt_func)
	}

	for _, f := range tgt_func {
		AnalyseFunction(f, tgt_pkg)
	}

	// Filter out not Applicable Functions
	// Extract AST of Target Function
	// Extract All ifs from the Function

	return
}

func Initialize() *Context {
	flag.Parse()
	if *flagArgumentType == "" {
		*flagArgumentType = "testing.T"
	}
	if *flagFunc != "" {
		*flagFuncPrefix 	= ""
		*flagArgumentType 	= ""
	}

	c := new(Context)
	c.makeWorkdir()
	// TODO: Load the "*_test.go" files along with the general "*.go" files
	return c
	// cfg := &packages.Config{Mode: packages.LoadAllSyntax}
	// pkgs, err := packages.Load(cfg, "testing")
	// if err != nil {
	// 	panic(fmt.Sprintf("Package \"testing\" not found"))
	// }
	// if len(pkgs) != 1 {
	// 	paths := make([]string, len(pkgs))
	// 	for _, p := range pkgs {
	// 		paths = append(paths, p.PkgPath)
	// 	}
	// 	panic(fmt.Sprintf("Cannot build multiple packages, %q resolved to %v", "testing", strings.Join(paths, ", ")))
	// }
	// pkg := pkgs[0]
	// for name, def := range pkg.TypesInfo.Defs {
	// 	if name.Name != "T" {
	// 		continue
	// 	}
	// 	a := def.Type().(*types.Named)
	// 	fmt.Printf("%s\n", a.String())
	// }
}

func main() {
	_ = Initialize()
	ExtractFunctionFromPackage()
}

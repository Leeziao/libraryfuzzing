package main

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type FunctionExtractor struct {
	f        *ast.File
	p        *packages.Package
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

func (fe *FunctionExtractor) Visit(node ast.Node) ast.Visitor {
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

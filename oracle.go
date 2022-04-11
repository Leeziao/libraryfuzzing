package main

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"os"

	"golang.org/x/tools/go/packages"
)

type Oracle struct {
	f_ast *ast.FuncDecl
	pkg *packages.Package
	fset *token.FileSet
}

func (o *Oracle) Visit(node ast.Node) ast.Visitor {
	switch st := node.(type) {
	case *ast.IfStmt:
		printer.Fprint(os.Stdout, o.fset, st.Cond)
		fmt.Printf("\n")
		// ast.Print(nil, st.Cond)
		fmt.Printf("%s\n", st.Cond)
		ast.Inspect(st.Cond, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.Ident:
				fmt.Printf("\t%s\n", n.Name)
			case *ast.CallExpr:
				return false
			}
			return true
		})
	}
	return o
}

func (o *Oracle) AnalyseFunction() {
	f_ast := o.f_ast
	fmt.Printf("--- Analysing %q ---\n", f_ast.Name.Name)
	if *flagShowSource {
		fmt.Printf("--- Begin: Source of %q ---\n", f_ast.Name.Name)
		printer_conf := printer.Config{
							Mode: printer.UseSpaces,
							Tabwidth: 4}
		printer_conf.Fprint(os.Stdout, o.pkg.Fset, f_ast)
		fmt.Printf("\n--- End: Source of %q ---\n", f_ast.Name.Name)
	}

	ast.Walk(o, f_ast)
}

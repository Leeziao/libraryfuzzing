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
	}
	return o
}
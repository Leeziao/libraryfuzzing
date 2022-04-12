package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"os"

	"golang.org/x/tools/go/packages"
)

type Oracle struct {
	f_ast *ast.FuncDecl
	pkg *packages.Package
	fset *token.FileSet
	f_scope *types.Scope
}

func (o *Oracle) Visit(node ast.Node) ast.Visitor {
	switch st := node.(type) {
	case *ast.IfStmt:
		var buf bytes.Buffer
		buf.WriteString("Original Condition: ")
		printer.Fprint(&buf, o.fset, st.Cond)
		buf.WriteString("\n")
		PrintColoredText(buf.String(), RED)

		o.InspectCondition(st.Cond)		// Exhibit variables' delete attribute
		expr, replaced := o.TraverseCondition(st.Cond) // Do the delete work

		if replaced {
			st.Cond = ast.NewIdent("false")
		} else {
			st.Cond = expr
		}

		buf.Reset()
		buf.WriteString("Modified Condition: ")
		printer.Fprint(&buf, o.fset, st.Cond)
		buf.WriteString("\n")
		PrintColoredText(buf.String(), GREEN)

		// ast.Print(nil, st.Cond)
		// ast.Walk(o, st.Body)
		// fmt.Printf("%s\n", st.Cond)
	}
	return o
}

// ret1 := modified expr
// ret2 := if current node is modified to be always accepted
func (o *Oracle) TraverseCondition(expr ast.Expr) (ast.Expr, bool) {
	switch expr := expr.(type) {
	case *ast.Ident:
		is_kept := o.IsExprKept(expr)
		if !is_kept {
			return ast.NewIdent("true"), true
		}
		return expr, false
	case *ast.BasicLit:
		return ast.NewIdent("true"), true
		return expr, false
	case *ast.BinaryExpr:
		x_expr, x_replaced := o.TraverseCondition(expr.X)
		y_expr, y_replaced := o.TraverseCondition(expr.Y)

		if x_replaced && y_replaced {
			return ast.NewIdent("true"), true
		}
		if (!x_replaced) && (!y_replaced) {
			expr.X = x_expr
			expr.Y = y_expr
			return expr, false
		}

		switch expr.Op {
		case token.LOR:
			if x_replaced {
				return y_expr, false
			} else {
				return x_expr, false
			}
		case token.LAND:
			if x_replaced {
				return y_expr, false
			} else {
				return x_expr, false
			}
		case token.LSS: fallthrough
		case token.GTR: fallthrough
		case token.LEQ: fallthrough
		case token.GEQ: fallthrough
		case token.EQL: fallthrough
		case token.NEQ:
			return ast.NewIdent("true"), true
		}
	case *ast.UnaryExpr:
		if expr.Op != token.NOT {
			return expr, false
		}
		sub_expr, sub_replaced := o.TraverseCondition(expr.X)
		if sub_replaced {
			return ast.NewIdent("true"), true
		}
		if sub_ident, ok := sub_expr.(*ast.Ident); ok {
			if sub_ident.Name == "true" {
				return ast.NewIdent("false"), false
			} else if sub_ident.Name == "false" {
				return ast.NewIdent("true"), false
			} else {
				return expr, false
			}
		}
	case *ast.StarExpr:
		_, sub_replaced := o.TraverseCondition(expr.X)
		if sub_replaced {
			return expr, true // TODO: StarExpr
		}
	case *ast.ParenExpr:
		sub_expr, sub_replaced := o.TraverseCondition(expr.X)
		if sub_replaced {
			return ast.NewIdent("true"), true
		}
		return sub_expr, false
	// default:
	// 	panic("Other undefined cases")
	}
	return expr, false
}

func (o *Oracle) InspectCondition(cond ast.Expr) {
	ast.Inspect(cond, func(nn ast.Node) bool {
		n, ok := nn.(ast.Expr)
		if !ok {
			return true
		}
		is_kept := o.IsExprKept(n)

		switch n := n.(type) {
		case *ast.SelectorExpr:
	 		var buf bytes.Buffer
			printer.Fprint(&buf, o.fset, n.X)
			buf.WriteString(fmt.Sprintf(".%s", n.Sel.Name))
			
			PrintColoredText(fmt.Sprintf("\t%s (kept=%v)\n", buf.String(), is_kept), YELLOW)
			return false
		case *ast.BasicLit:
			PrintColoredText(fmt.Sprintf("\t%s (kept=%v)\n", n.Value, is_kept), YELLOW)
		case *ast.Ident:
			is_kept := o.IsExprKept(n)
			PrintColoredText(fmt.Sprintf("\t%s (kept=%v)\n", n.Name, is_kept), YELLOW)
		}
		return true
	})
}

func (o *Oracle) IsExprKept(expr ast.Expr) bool {
	switch expr := expr.(type) {
	case *ast.SelectorExpr:
		return true
	case *ast.Ident:
		pos := expr.NamePos
		pkg_scope := o.pkg.Types.Scope()
		inner_scope := pkg_scope.Innermost(pos)

		define_s, obj := inner_scope.LookupParent(expr.Name, pos)
		// Keep if is a local variable
		if o.f_scope.Contains(obj.Pos()) {
			return true
		}
		// Keep if is a universe variable
		for _, name := range types.Universe.Names() {
			if name == expr.Name {
				return true
			}
		}
		return false
		_ = obj
		_ = define_s
	case *ast.BasicLit:
		return false
	}
	return true
}

func (o *Oracle) AnalyseFunction() {
	f_ast := o.f_ast
	fmt.Printf("--- Analysing %q ---\n", f_ast.Name.Name)
	if *flagShowSource {
		fmt.Printf("--- (Before Modify) Begin: Source of %q ---\n", f_ast.Name.Name)
		printer_conf := printer.Config{
							Mode: printer.UseSpaces,
							Tabwidth: 4}
		printer_conf.Fprint(os.Stdout, o.pkg.Fset, f_ast)
		fmt.Printf("\n--- (Before Modify) End: Source of %q ---\n", f_ast.Name.Name)
	}

	// TODO: The Workflow of cond modifying
	// TODO: Now only binary operands are considered
	// 1. Traverse in DFS order
	// 2. ...
	ast.Walk(o, f_ast)

	if *flagShowSource {
		fmt.Printf("--- (After Modify) Begin: Source of %q ---\n", f_ast.Name.Name)
		printer_conf := printer.Config{
							Mode: printer.UseSpaces,
							Tabwidth: 4}
		printer_conf.Fprint(os.Stdout, o.pkg.Fset, f_ast)
		fmt.Printf("\n--- (After Modify) End: Source of %q ---\n", f_ast.Name.Name)
	}
}

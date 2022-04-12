package main

import (
	"flag"
	"fmt"
	"go/types"
)

var (
	flagFunc			= flag.String("f", "", "Function to analyse")
	flagFuncPrefix 		= flag.String("p", "", "Prefix to keep function")
	flagArgumentType 	= flag.String("arg", "", "Type of the argument of target functions")
	flagShowSource 		= flag.Bool("s", false, "Show source code of analysed function")
)

func OtherConditions() {}

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
	c.MakeWorkdir()
	// TODO: Load the "*_test.go" files along with the general "*.go" files
	return c
}

func main() {
	// Filter out not Applicable Functions
	// Extract AST of Target Function
	// Extract All Oracles from the Function
	// Analyse the Oracles

	context := Initialize()
	tgt_func := context.ExtractFunctionFromPackage()

	for _, f := range tgt_func {
		tgt_pkg := context.tgt_pkg
		tgt_info := tgt_pkg.TypesInfo
		f_obj := tgt_info.Defs[f.Name]
		f_func, ok := f_obj.(*types.Func)

		if !ok {
			panic(fmt.Sprintf("Function Object of %s not found", f.Name.Name))
		}


		oracle := Oracle{f_ast: f,
						 fset: tgt_pkg.Fset,
						 pkg: tgt_pkg,
						 f_scope: f_func.Scope(),}
		oracle.AnalyseFunction()
		break
	}
}

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Context struct {
	workdir 	string 				// TODO: a temporary workdir
	pkgs 		[]*packages.Package // TODO: packages that tgt_pkg rely on
	tgt_pkg 	*packages.Package
}

func (c *Context) MakeWorkdir() {
	var err error
	c.workdir, err = ioutil.TempDir("", "libraryfuzzing")
	if err != nil {
		panic(fmt.Sprintf("Failed to create tmp dir: %v", err))
	}
}

func (c *Context) cleanup() {
	os.RemoveAll(c.workdir)
}

func (c *Context) ExtractFunctionFromPackage() []*ast.FuncDecl {
	c.tgt_pkg = c.GetPackage()
	var tgt_func []*ast.FuncDecl

	for i, fullname := range c.tgt_pkg.CompiledGoFiles {
		fname := filepath.Base(fullname)
		if !strings.HasSuffix(fname, ".go") {
			// This is a cgo-generated file.
			// Currently does not work.
			// We copied the original Go file as part of copyPackageRewrite,
			// so we can just skip this one.
			// See https://golang.org/issue/30479.
			continue
		}
		f := c.tgt_pkg.Syntax[i]
		ExtractFunctionFromFile(f, c.tgt_pkg, &tgt_func)
	}
	return tgt_func
}

func (c *Context) GetPackage() *packages.Package {
	if flag.NArg() > 1 {
		panic("Should not receive more than one argument")
	}
	pkg_name := "."
	if flag.NArg() == 1 {
		pkg_name = flag.Arg(0)
	}
	fmt.Printf("Loading package %q\n", pkg_name)

	cfg := &packages.Config{
				Mode: packages.LoadAllSyntax,
				Tests: true}
	pkgs, err := packages.Load(cfg, pkg_name)
	if err != nil {
		panic(fmt.Sprintf("Package %v not found", flag.Args()))
	}

	if len(pkgs) != 4 {
		paths := make([]string, len(pkgs))
		for _, p := range pkgs {
			paths = append(paths, p.PkgPath)
		}
		panic(fmt.Sprintf("Cannot build multiple packages, %q resolved to %v", pkg_name, strings.Join(paths, ", ")))
	}

	pkg := pkgs[1]
	for _, name := range pkg.CompiledGoFiles {
		fmt.Printf("\t- %s\n", name)
	}
	return pkg
}

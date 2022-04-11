package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"golang.org/x/tools/go/packages"
)

type Context struct {
	workdir string
	pkgs []*packages.Package
}

func (c *Context) makeWorkdir() {
	var err error
	c.workdir, err = ioutil.TempDir("", "libraryfuzzing")
	if err != nil {
		panic(fmt.Sprintf("Failed to create tmp dir: %v", err))
	}
}

func (c *Context) cleanup() {
	os.RemoveAll(c.workdir)
}

func TestOK(t *testing.T) {
	fmt.Println("TestOK is called")
	a := 19
	if b := a + 3; b > 10 {
		fmt.Println("Good")
	} else {
		fmt.Println("Bad")
	}
}
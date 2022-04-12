package main

import (
	"fmt"
	"testing"
)

// binary_op  = "||" | "&&" | rel_op | add_op | mul_op .
// rel_op     = "==" | "!=" | "<" | "<=" | ">" | ">=" .
// add_op     = "+" | "-" | "|" | "^" .
// mul_op     = "*" | "/" | "%" | "<<" | ">>" | "&" | "&^" .

// unary_op   = "+" | "-" | "!" | "^" | "*" | "&" | "<-" .

var g_V1 int = 9
var g_V2 *int

// a > 10 -> a > 10
// g_a > 10 -> true
// a > 10 && g_a > 10 -> a > 10 && true
// a > 10 || g_a > 10 -> a > 10 || false
// !(true) -> !(false)
// !(false) -> !(true)
// !(a > 3) -> (a > 3) != true

// if b := a + 3; !(g_V1 == 8 || g_V2 == nil) || b > 10 && a < 9 && *flagShowSource || a == b {
func Test1(t *testing.T) {
	a := 19
	if b := a + 3; !(a == b && a == 20) {
		_ = b
		fmt.Println("Good")
	} else {
		fmt.Println("Bad")
	}
}
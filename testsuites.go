package main

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	fmt.Println("Test1 is called")
	a := 19
	if b := a + 3; b > 10 {
		fmt.Println("Good")
	} else {
		fmt.Println("Bad")
	}
}
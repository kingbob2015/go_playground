package main

import (
	"fmt"
)

var packageVariable string

func init() {
	packageVariable = "hello"
	//Since short declarations are only local, this declares a new variable, doesnt assign the global
	//Short declarations and REQUIRED to declare a NEW variable. If all the variables in a short declaration
	//already exist then it is not a valid statement
	packageVariable := 60
	fmt.Println(packageVariable)
}

func mainPS() {
	a, b := 0, "test"
	pointerFunc(&b)
	fmt.Printf("%d %s package variable: %s\n", a, b, packageVariable)
	p := new(int)
	*p = 200
	fmt.Println(*p)
}

func pointerFunc(s *string) {
	fmt.Println(s)
	*s += "bobwashere"
}

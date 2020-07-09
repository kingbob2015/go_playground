package main

import (
	"fmt"
	"os"
)

func mainEcho() {
	for i, arg := range os.Args[1:] {
		fmt.Print(i)
		fmt.Println(" " + arg)
	}
}

// strings.Join is a much more efficient opperation that building up a string with a separator ourselves
// func main() {
// 	fmt.Println(os.Args[0])
// 	fmt.Println(strings.Join(os.Args[1:], " "))
// }

package main

import (
	"bytes"
	"fmt"
)

func mainStrings() {
	s := "Hello World"
	b := []byte(s)
	fmt.Println(s)
	fmt.Println(b)
	//Cant do b[0] = "Y" because we need it to be a rune
	//In UTF-8 the first bit for an ASCII character is a 0 to show it is ASCII for next 7 bits
	b[0] = 'Y'
	s = string(b)
	fmt.Println(s)
	for i := 0; i < len(b); i++ {
		fmt.Println(b[i])
	}
	//Slices are dynamically sized so you could do the additions with the byte slice but the Buffer makes
	//it simpler

	fmt.Println(bytesComma("123456789"))
}

func bytesComma(s string) string {
	//Can write strings and byte and []byte and the like to a Buffer
	var buf bytes.Buffer
	leadingDigits := len(s) % 3
	buf.WriteString(s[:leadingDigits])
	for i := leadingDigits; i < len(s); i += 3 {
		if i != 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(s[i : i+3])
	}
	return buf.String()
}

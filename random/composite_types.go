package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	arr := [...]int{1, 2, 3, 4, 5}
	reverseWithSlice(arr[:])
	fmt.Println(arr)
	reverseWithArrPointer(&arr)
	fmt.Println(arr)
	rotate(arr[:], 2)
	fmt.Println(arr)

	file := os.Args[1]
	if len(file) != 0 {
		f, err := os.Open(file)
		if err != nil {

		} else {
			wordFreq(f)
		}
	}
}

func wordFreq(f *os.File) {
	input := bufio.NewScanner(f)
	counts := map[string]int{}
	input.Split(bufio.ScanWords)
	for input.Scan() {
		counts[input.Text()]++
	}
	fmt.Print(counts)
}

func rotate(s []int, x int) {
	firstPart := s[x:]
	secondPart := s[:x]
	for i, v := range firstPart {
		s[i] = v
	}
	for i, v := range secondPart {
		s[i+len(firstPart)] = v
	}
}

func reverseWithSlice(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func reverseWithArrPointer(arr *[5]int) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

// func modArrPointer(ptr *[32]byte) {
// 	fmt.Println(ptr)
// 	//Both these work. Range must auto get value at pointer?
// 	for _, v := range ptr {
// 		fmt.Println(v)
// 	}
// 	for _, v := range *ptr {
// 		fmt.Println(v)
// 	}
// }

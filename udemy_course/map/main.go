package main

import "fmt"

func main() {
	// Difference between these declarations is the first does not allocate memory and points to nil. Second, make, allocates memory and points to a map with 0 elements.
	// var colors map[string]string
	// colors := make(map[string]string)

	colors := map[string]string{
		"red":   "#ff0000",
		"green": "#4b756d",
	}

	colors["white"] = "#ffffff"
	// delete(colors, "white")

	printMap(colors)
}

func printMap(m map[string]string) {
	for k, v := range m {
		fmt.Printf("Hex code for %v is %v\n", k, v)
	}
}

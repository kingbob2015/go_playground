package main

import "fmt"

type shape interface {
	getArea() float64
}

type triangle struct {
	height float64
	base   float64
}
type square struct {
	sideLength float64
}

func main() {
	t := &triangle{
		height: 11.0,
		base:   13.0,
	}

	s := &square{
		sideLength: 25.0,
	}

	printArea(t)
	printArea(s)
}

func (t *triangle) getArea() float64 {
	return t.base * t.height * 0.5
}

func (s *square) getArea() float64 {
	return s.sideLength * s.sideLength
}

func printArea(s shape) {
	fmt.Printf("The area of the shape is %v\n", s.getArea())
}

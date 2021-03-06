package main

import "fmt"

type bot interface {
	getGreeting() string
	setTest(string)
}

type englishBot struct{ test string }
type spanishBot struct{ test string }

func main() {
	eb := englishBot{}
	sb := spanishBot{}

	//Since the methods are using pointer types then when b bot gets called, the interface is satisfied by *eb and *sb instead of just eb and sb
	//If we wanted to work with just eb and just sb as well would need other receiver functions for that
	//Also we can't do overloading for parameter functions so receiver functions are needed for a type of overloading
	printGreeting(&eb)
	printGreeting(&sb)
	//Go automatically adds the &sb to call the pointer receiver function on a value
	sb.setTest("Whatsup")
	setTestCall(&eb)
	fmt.Println(eb)
	fmt.Println(sb)
}

func printGreeting(b bot) {
	fmt.Println(b.getGreeting())
}

func setTestCall(b bot) {
	b.setTest("Test")
}

func (eb *englishBot) setTest(s string) {
	eb.test = s
}

//If a receiver function is a pointer, it only works on pointers
func (sb *spanishBot) setTest(s string) {
	sb.test = s
}

//If a receiver function is a value it works on values and pointers
func (englishBot) getGreeting() string {
	return "Hello"
}

func (spanishBot) getGreeting() string {
	return "Hola"
}

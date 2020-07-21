package main

func main() {
	cards := newDeckFromFile("hand.txt")
	cards.shuffle()
	cards.print()
}

package main

import (
	"os"
	"testing"
)

func TestNewDeck(t *testing.T) {
	d := newDeck()
	if len(d) != 52 {
		t.Errorf("Expected deck length of 52, but got %v", len(d))
	}

	if d[0] != "Ace of Spades" {
		t.Errorf("Expected Ace of Spades, but got %v", d[0])
	}

	if d[len(d)-1] != "King of Clubs" {
		t.Errorf("Expected King of Clubs, but got %v", d[len(d)-1])
	}
}

func TestNewDeckFromFileAndSaveToFile(t *testing.T) {
	os.Remove("_decktesting")

	d := newDeck()
	d.saveToFile("_decktesting")

	loadedDeck := newDeckFromFile("_decktesting")

	if len(loadedDeck) != 52 {
		t.Errorf("Expected deck length of 52, but got %v", len(loadedDeck))
	}

	if loadedDeck[0] != "Ace of Spades" {
		t.Errorf("Expected Ace of Spades, but got %v", loadedDeck[0])
	}

	if loadedDeck[len(d)-1] != "King of Clubs" {
		t.Errorf("Expected King of Clubs, but got %v", loadedDeck[len(d)-1])
	}

	os.Remove("_decktesting")

}

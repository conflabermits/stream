package main

import (
	"testing"
)

func TestAlphabetizeShort(t *testing.T) {
	// Alphabetize the words in a short string
	input := "Delta Iota Eta Epsilon Alpha zeta gamma beta theta"
	expected := "Alpha beta Delta Epsilon Eta gamma Iota theta zeta"
	output := alphabetize(input)
	if output != expected {
		t.Errorf("Expected:\n%q\n\nGot:\n%q", expected, output)
	}
}

func TestAlphabetizeLong(t *testing.T) {
	// Alphabetize the words in a longer quote
	input := "Random quote -- Knowing your own darkness is the best method for dealing with the darknesses of other people... in bed."
	expected := "-- bed. best darkness darknesses dealing for in is Knowing method of other own people... quote Random the the with your"
	output := alphabetize(input)
	if output != expected {
		t.Errorf("Expected:\n%q\n\nGot:\n%q", expected, output)
	}
}

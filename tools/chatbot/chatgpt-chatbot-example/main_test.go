package main

import (
	"testing"
)

func TestAlphabetize(t *testing.T) {
	// Test case 1: Alphabetize the words in a short string
	input1 := "Delta Iota Eta Epsilon Alpha zeta gamma beta theta"
	expected1 := "Alpha beta Delta Epsilon Eta gamma Iota theta zeta"
	output1 := alphabetize(input1)
	if output1 != expected1 {
		t.Errorf("Expected:\n%q\n\nGot:\n%q", expected1, output1)
	}

	// Test case 2: Alphabetize the words in a longer quote
	input2 := "Random quote -- Knowing your own darkness is the best method for dealing with the darknesses of other people... in bed."
	expected2 := "-- bed. best darkness darknesses dealing for in is Knowing method of other own people... quote Random the the with your"
	output2 := alphabetize(input2)
	if output2 != expected2 {
		t.Errorf("Expected:\n%q\n\nGot:\n%q", expected2, output2)
	}
}

package main

import (
	"testing"
)

func TestAlphabetize(t *testing.T) {
	// Test case 1: Alphabetize the words in a longer quote
	input1 := "Random quote -- Knowing your own darkness is the best method for dealing with the darknesses of other people... in bed."
	expected1 := "-- bed. best darkness darknesses dealing for in is Knowing method of other own people... quote Random the the with your"
	output1 := alphabetize(input1)
	if output1 != expected1 {
		t.Errorf("Expected:\n%q\n\nGot:\n%q", expected1, output1)
	}
}

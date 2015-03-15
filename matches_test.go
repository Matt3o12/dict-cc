package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchTerm(t *testing.T) {
	term := "Hello world. I'm a test containing umlauts: äëïöü."
	expected := "Hello+world.+I%27m+a+test+containing+umlauts%3A+%C3%A4%C3%AB" +
		"%C3%AF%C3%B6%C3%BC."
	assert.Equal(t, expected, escapeSearchTerm(term))
}

func TestGetLookupURLForSearchTerm(t *testing.T) {
	pair := makePair("English", "German")
	pair.Second.Abbrev = "de"
	term := "Hello World"

	expected := "http://ende.dict.cc/?s=Hello+World"
	assert.Equal(t, expected, GetLookupURLForSearchTerm(pair, term))
}

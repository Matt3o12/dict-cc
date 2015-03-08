package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLanguagePairByString(t *testing.T) {
	// expected = LanguagePair{Language{"German"}
}

func TestLanguagePairSame(t *testing.T) {
	pair1 := LanguagePair{"German", "English"}
	pair2 := LanguagePair{"German", "English"}
	assert.True(t, pair1.Same(pair2))
	assert.True(t, pair2.Same(pair1))
}

func TestLanguagePairSameSwapped(t *testing.T) {
	pair1 := LanguagePair{"German", "English"}
	pair2 := LanguagePair{"English", "German"}
	assert.True(t, pair1.Same(pair2))
	assert.True(t, pair2.Same(pair1))
}

func TestLanguagePairSameFail(t *testing.T) {
	pair1 := LanguagePair{"German", "English"}
	pair2 := LanguagePair{"German", "Spanish"}
	assert.False(t, pair1.Same(pair2))
	assert.False(t, pair2.Same(pair1))
}

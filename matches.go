package main

import (
	"fmt"
	"net/url"
	"strings"
)

// LanguageResult represents the result when you looked up
// a word.
type LanguageResult struct {
	Needle  string
	Matches []string
	Pair    LanguagePair
}

// Path used for looking up words.
const DictCCWorldLookupPatch = "http://%v%v.dict.cc/?s=%v"

func escapeSearchTerm(searchTerm string) string {
	searchTerm = url.QueryEscape(searchTerm)

	return strings.Replace(searchTerm, "%2B", "+", -1)
}

// GetLookupURLForSearchTerm generates the urls used for looking up
// the search term `searchTerm` in the given language
func GetLookupURLForSearchTerm(pair LanguagePair, searchTerm string) string {
	searchTerm = escapeSearchTerm(searchTerm)
	a1, a2 := pair.First.Abbrev, pair.Second.Abbrev

	return fmt.Sprintf(DictCCWorldLookupPatch, a1, a2, searchTerm)
}

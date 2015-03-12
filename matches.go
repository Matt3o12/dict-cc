package main

import "net/http"

// LanguageResult represents the result when you looked up
// a word.
type LanguageResult struct {
	Needle  string
	Matches []string
	Pair    LanguagePair
}

// FindResults finds all results on the result page.
func FindResults(response *http.Response) (*LanguageResult, error) {

	return nil, nil
}

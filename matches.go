package main

import (
	"fmt"
	"net/http"
)

// LanguageResult represents the result when you looked up
// a word.
type LanguageResult struct {
	Needle  string
	Matches []string
	Pair    LanguagePair
}

const DictCCWorldLookupPatch = "http://%v.dict.cc/?s=%v"

// FindResults finds all results on the result page.
func FindResults(response *http.Response) (*LanguageResult, error) {
	_ = fmt.Sprintf(DictCCWorldLookupPatch)

	return nil, nil
}

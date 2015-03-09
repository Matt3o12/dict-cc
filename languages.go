package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// A Language which includes the localized (English) and native
// version of the language
type Language string

// A LanguagePair consists of two langauges.
type LanguagePair struct {
	first, second Language
}

// Same checks whether the current language pair equals the other one.
// it swaps the first and second language if necessary.
func (l LanguagePair) Same(other LanguagePair) bool {
	first1, second1 := l.first, l.second
	first2, second2 := other.first, other.second
	if first1 != first2 {
		// Swapping first1 and first2
		first2, first1 = first1, first2
	}

	return (first1 == first2 && second1 == second2) || (first1 == second1 && first2 == second2)
}

func (l LanguagePair) String() string {
	return fmt.Sprintf("{%v - %v}", l.first, l.second)
}

// GetLanguagesFromRemote returns all avaiable
// languages found in the response.
func GetLanguagesFromRemote(response *http.Response) ([]LanguagePair, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}

	var pairs []LanguagePair
	var langErr *error

	doc.Find(allAvaiableLangsCSSPath).Each(
		func(i int, s *goquery.Selection) {
			lang, err := getLanguagePairFromString(s.Text())
			if err != nil {
				langErr = &err
				return
			}

			pairs = append(pairs, *lang)
		})

	if langErr != nil {
		return nil, *langErr
	}

	return pairs, nil
}

func getLanguagePairFromString(s string) (*LanguagePair, error) {
	split := strings.Split(s, "–") // that's not a normal hython ('-').
	if len(split) != 2 {
		return nil, fmt.Errorf("Unkown language format: '%v'. "+
			"Are you using the latest version?", s)
	}

	first := Language(strings.TrimSpace(split[0]))
	second := Language(strings.TrimSpace(split[1]))
	pair := &LanguagePair{first, second}
	return pair, nil
}

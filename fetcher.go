package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	// DictBaseURL The base url of dict.cc
	DictBaseURL = "http://dict.cc/"

	// AllLangaugesGet URL where all available langauge pairs can be found.
	AllLangaugesGet = "http://browse.dict.cc/"

	allAvaiableLangsCSSPath = "#maincontent form[name='langbarchooser'] " +
		"table td a"
)

// A Language which includes the localized (English) and native
// version of the language
type Language struct {
	native string
}

// A LanguagePair consists of two langauges.
type LanguagePair struct {
	first, second Language
}

func (l Language) String() string {
	return l.native
}

func (l LanguagePair) String() string {
	return fmt.Sprintf("{%v - %v}", l.first, l.second)
}

// GetLanguages returns all avaiable
// languages found on AllLangaugesGet.
func GetLanguages() ([]LanguagePair, error) {
	doc, err := goquery.NewDocument(AllLangaugesGet)

	if err != nil {
		return nil, err
	}

	var pairs []LanguagePair
	var langErr *error

	doc.Find(allAvaiableLangsCSSPath).Each(
		func(i int, s *goquery.Selection) {
			lang, err := getLanguagePairByString(s.Text())
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

func getLanguagePairByString(s string) (*LanguagePair, error) {
	split := strings.Split(s, "–") // that's not the normal hython.
	if len(split) != 2 {
		fmt.Println(split)
		return nil, fmt.Errorf("Unkown language format: '%v'. "+
			"Are you using the latest version?", s)
	}

	first := Language{strings.TrimSpace(split[0])}
	second := Language{strings.TrimSpace(split[1])}

	pair := &LanguagePair{first, second}
	return pair, nil
}

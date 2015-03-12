package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// OutdatedLanguageFileError is used to indecate that the
// language file needs updating
var ErrOutdatedLanguageFile = errors.New(
	"The language pair file needs updating.")

const (
	// AllLangaugesGet URL where all available langauge pairs can be found.
	AllLangaugesGet = "http://browse.dict.cc/"

	allAvaiableLangsCSSPath = "#maincontent form[name='langbarchooser'] " +
		"table td a"

	languageFileVersion = 1
)

// A Language which includes the localized (English) and native
// version of the language
type Language struct {
	Name   string
	Abbrev string
}

func (l Language) String() string {
	return l.Name
}

// A LanguagePair consists of two langauges.
type LanguagePair struct {
	First, Second Language
}

type languageFileFormat struct {
	Version int
	Pairs   []LanguagePair
}

// Same checks whether the current language pair equals the other one.
// it swaps the first and second language if necessary.
func (l LanguagePair) Same(other LanguagePair) bool {
	first1, second1 := l.First, l.Second
	first2, second2 := other.First, other.Second
	if first1 != first2 {
		// Swapping first1 and first2
		first2, first1 = first1, first2
	}

	return (first1 == first2 && second1 == second2) || (first1 == second1 && first2 == second2)
}

func (l LanguagePair) String() string {
	return fmt.Sprintf("{%v - %v}", l.First, l.Second)
}

// LoadLanguagesFromDisk loads all languages from the disk.
func LoadLanguagesFromDisk(reader io.Reader) ([]LanguagePair, error) {
	decoder := json.NewDecoder(reader)
	var data languageFileFormat
	// Fixme: data.Version branch is not tested.
	if err := decoder.Decode(&data); err != nil && data.Version != languageFileVersion {
		return nil, ErrOutdatedLanguageFile
	}

	return data.Pairs, nil
}

// SaveLanguagesToDisk saves the language to disk (or the given writer).
func SaveLanguagesToDisk(langs []LanguagePair, writer io.Writer) error {
	data := languageFileFormat{Version: languageFileVersion, Pairs: langs}
	encoder := json.NewEncoder(writer)
	return encoder.Encode(data)
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
	firstS, secondS, err := splitLangauge(s)
	if err != nil {
		return nil, err
	}

	first := Language{Name: firstS, Abbrev: firstS[:2]}
	second := Language{Name: secondS, Abbrev: secondS[:2]}
	pair := &LanguagePair{first, second}

	return pair, nil
}

func splitLangauge(s string) (first, second string, err error) {
	split := strings.Split(s, "–") // that's not a normal hython ('-').
	if len(split) != 2 {
		return "", "", fmt.Errorf("Unkown language format: '%v'. "+
			"Are you using the latest version?", s)
	}

	first = strings.TrimSpace(split[0])
	second = strings.TrimSpace(split[1])

	return
}

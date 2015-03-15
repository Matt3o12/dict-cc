package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// OutdatedLanguageFileError is used to indecate that the
// language file needs updating
var ErrOutdatedLanguageFile = errors.New(
	"The language pair file needs updating.")

const (
	// AllLangaugesGet URL where all available langauge pairs can be found.
	AllLangaugesGet = "http://www.dict.cc/"

	allLangsCSSBasePath = "#maincontent table:nth-child(7) tr:nth-child(2)"
	englishPairCSSPath  = "td:nth-child(1) a"
	germanPairCSSPath   = "td:nth-child(2) a"

	abbrevFinderRegex    = `^http\:\/\/([a-zA-Z]{2})([a-zA-Z]{2})\.dict\.cc`
	abbrevSpecialCaseWWW = `^http\:\/\/www\.dict\.cc`

	languageFileVersion = 2

	parsingErrorStr = "There was an error " +
		"parsing the page. Are you using the latest version? (Error-ID: %v)"
)

var (
	germanLang  = Language{"Deutsch", "de"}
	englishLang = Language{"English", "en"}
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

func getPairsFromSelectors(firstLang Language, results []LanguagePair) func(int, *goquery.Selection) {
	return func(n int, selection *goquery.Selection) {

	}
}

// makeError creates a new error using message and returns the reference. That is a
// workaround for
// https://groups.google.com/forum/#!topic/golang-nuts/OGshFH1Z-Mc
func makeError(message string) *error {
	err := errors.New(message)
	return &err
}

func makeParsingError(id int) *error {
	msg := fmt.Sprintf(parsingErrorStr, id)
	return makeError(msg)
}

// GetLanguagesFromRemote returns all avaiable
// languages found in the response.
func GetLanguagesFromRemote(response *http.Response) ([]LanguagePair, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}

	base := doc.Find(allLangsCSSBasePath)
	pairPathes := map[string]Language{
		englishPairCSSPath: englishLang,
		germanPairCSSPath:  germanLang,
	}

	var pairs []LanguagePair
	var langErr *error
	urlRegex := regexp.MustCompile(abbrevFinderRegex)
	urlRegexWWW := regexp.MustCompile(abbrevSpecialCaseWWW)

	pairs = append(pairs, LanguagePair{germanLang, englishLang})
	for path, firstLang := range pairPathes {
		results := base.Find(path)
		results.Each(func(i int, s *goquery.Selection) {
			if err != nil {
				return
			}

			link, exists := s.Attr("href")
			if !exists {
				langErr = makeParsingError(1)
				return
			}

			matches := urlRegex.FindStringSubmatch(link)
			if len(matches) == 3 {
				abbrev, err := getAbbrev(firstLang, matches[1:])
				if err != nil {
					panic(err)
				}

				secondLange := Language{strings.TrimSpace(s.Text()), abbrev}
				pairs = append(pairs, LanguagePair{firstLang, secondLange})
			} else if !urlRegexWWW.MatchString(link) {
				langErr = makeParsingError(2)
			}
		})
	}

	if langErr != nil {
		return nil, *langErr
	}

	return pairs, nil
}

func getAbbrev(firstLang Language, matches []string) (string, error) {
	if len(matches) != 2 {
		return "", fmt.Errorf("Too many/few matches given (matches: %v).", matches)
	}

	first, second := matches[0], matches[1]
	if first == firstLang.Abbrev {
		return second, nil
	}

	return first, nil
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

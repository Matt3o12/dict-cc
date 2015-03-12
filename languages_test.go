package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"

	"github.com/matt3o12/dict-cc/plusTesting"
	"github.com/stretchr/testify/assert"
)

func makePair(first, second string) (pair LanguagePair) {
	firstL := Language{first, strings.ToLower(first[:2])}
	secondL := Language{second, strings.ToLower(second[:2])}

	pair = LanguagePair{firstL, secondL}
	return
}

func TestGetLanguagePairFromString(t *testing.T) {
	result, err := getLanguagePairFromString("German – English")
	assert.Nil(t, err)

	l1, l2 := Language{"German", "Ge"}, Language{"English", "En"}
	expected := LanguagePair{l1, l2}
	assert.True(t, result.Same(expected))
}

func TestGetLanguagePairFromStringFail(t *testing.T) {
	result, err := getLanguagePairFromString("German | English")
	assert.Nil(t, result)
	if assert.NotNil(t, err) {
		assert.EqualError(t, err, "Unkown language format: '"+
			"German | English'. Are you using the latest version?")
	}
}

func TestLanguagePairString(t *testing.T) {
	pair := makePair("German", "English")
	assert.Equal(t, "{German - English}", pair.String())
}

func TestLanguagePairSame(t *testing.T) {
	pair1 := makePair("German", "English")
	pair2 := makePair("German", "English")
	assert.True(t, pair1.Same(pair2))
	assert.True(t, pair2.Same(pair1))
}

func TestLanguagePairSameSwapped(t *testing.T) {
	pair1 := makePair("German", "English")
	pair2 := makePair("English", "German")
	assert.True(t, pair1.Same(pair2))
	assert.True(t, pair2.Same(pair1))
}

func TestLanguagePairSameFail(t *testing.T) {
	pair1 := makePair("German", "English")
	pair2 := makePair("German", "Spanish")
	assert.False(t, pair1.Same(pair2))
	assert.False(t, pair2.Same(pair1))
}

func fatalOnNil(t *testing.T, msg string, err error) {
	if err != nil {
		t.Fatal(msg, err)
	}
}

func TestGetLanguagesFromRemote(t *testing.T) {
	content, err := ioutil.ReadFile("resources/browse.dict.cc.html")
	fatalOnNil(t, "Error reading file.", err)

	uri, err := url.Parse(AllLangaugesGet)
	fatalOnNil(t, "Unexpected error parsing url", err)

	response := plusTesting.GetHTTPResponse(200, uri, string(content))
	langs, err := GetLanguagesFromRemote(response)
	assert.Nil(t, err)
	assert.Equal(t, 50, len(langs))
}

func TestGetLanguagesFromRemoteNilResponse(t *testing.T) {
	result, err := GetLanguagesFromRemote(nil)
	assert.Nil(t, result)
	assert.EqualError(t, err, "Response is nil pointer")
}

func TestGetLanguagesFromRemoteOutdatedFormat(t *testing.T) {
	content, err := ioutil.ReadFile("resources/browse-error.dict.cc.html")
	fatalOnNil(t, "Error reading file.", err)

	uri, err := url.Parse(AllLangaugesGet)
	fatalOnNil(t, "Unexpected error parsing url", err)

	response := plusTesting.GetHTTPResponse(200, uri, string(content))
	langs, err := GetLanguagesFromRemote(response)
	assert.Nil(t, langs)
	assert.EqualError(t, err, "Unkown language format: 'English | "+
		"Slovak'. Are you using the latest version?")
}

func loadLangs(rawData string) ([]LanguagePair, error) {
	return LoadLanguagesFromDisk(strings.NewReader(rawData))
}

func TestLoadLanguagesFromDisk(t *testing.T) {
	data := `{"Version":1,"Pairs":[{"First":{"Name":"German123","Abbrev":` +
		`"ge"},"Second":{"Name":"English321","Abbrev":"en"}},{"First":` +
		`{"Name":"German","Abbrev":"ge"},"Second":{"Name":"Russian",` +
		`"Abbrev":"ru"}}]}`

	langs, err := loadLangs(data)
	assert.Nil(t, err)
	expected := []LanguagePair{
		makePair("German123", "English321"),
		makePair("German", "Russian"),
	}
	assert.Equal(t, langs, expected)
}

func TestLoadLanguageFromDiskOldVersion(t *testing.T) {
	data := `{"Version":-1,"Pairs":[{"First":"German123","Second":` +
		`"English321"},{"First":"German","Second":"Russian"}]}`
	langs, err := loadLangs(data)
	assert.Equal(t, err, ErrOutdatedLanguageFile)
	assert.Nil(t, langs)
}

func TestSaveLanguages(t *testing.T) {
	buf := new(bytes.Buffer)
	langs := []LanguagePair{
		makePair("German", "English"),
		makePair("German", "Spanish"),
	}

	err := SaveLanguagesToDisk(langs, buf)
	assert.Nil(t, err)
	expected := `{"Version":%v,"Pairs":[{"First":{"Name":"German","Abbrev":` +
		`"ge"},"Second":{"Name":"English","Abbrev":"en"}},{"First":` +
		`{"Name":"German","Abbrev":"ge"},"Second":{"Name":"Spanish",` +
		`"Abbrev":"sp"}}]}`
	expected = fmt.Sprintf(expected, languageFileVersion)
	assert.Equal(t, strings.TrimSpace(buf.String()), expected)
}

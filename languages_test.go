package main

import (
	"bytes"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"

	"github.com/matt3o12/dict-cc/httpTesting"
	"github.com/stretchr/testify/assert"
)

func TestGetLanguagePairFromString(t *testing.T) {
	result, err := getLanguagePairFromString("German – English")
	assert.Nil(t, err)

	expected := LanguagePair{"German", "English"}
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
	pair := LanguagePair{"German", "English"}
	assert.Equal(t, "{German - English}", pair.String())
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

	response := httpTesting.GetHTTPResponse(200, uri, string(content))
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

	response := httpTesting.GetHTTPResponse(200, uri, string(content))
	langs, err := GetLanguagesFromRemote(response)
	assert.Nil(t, langs)
	assert.EqualError(t, err, "Unkown language format: 'English | "+
		"Slovak'. Are you using the latest version?")
}

func TestUpdateLanguagesIntregration(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test")
	}

	updateLanguages()
}

// TOOD: Add support for versionizing the languages.
func TestLoadLanguagesFromDisk(t *testing.T) {
	data := `[{"First":"German123","Second":"English321"},` +
		`{"First":"German","Second":"Russian"}]`
	langs, err := LoadLanguagesFromDisk(strings.NewReader(data))
	assert.Nil(t, err)
	expected := []LanguagePair{
		LanguagePair{"German123", "English321"},
		LanguagePair{"German", "Russian"},
	}
	assert.Equal(t, langs, expected)
}

func TestSaveLanguages(t *testing.T) {
	buf := new(bytes.Buffer)
	langs := []LanguagePair{
		LanguagePair{"German", "English"},
		LanguagePair{"German", "Spanish"},
	}

	err := SaveLanguagesToDisk(langs, buf)
	assert.Nil(t, err)
	expected := `[{"First":"German","Second":"English"},` +
		`{"First":"German","Second":"Spanish"}]`
	assert.Equal(t, strings.TrimSpace(buf.String()), expected)
}

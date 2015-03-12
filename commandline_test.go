package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/matt3o12/dict-cc/plusTesting"
	"github.com/stretchr/testify/assert"
)

type logWriter struct {
	*testing.T
}

func (writer logWriter) Write(p []byte) (n int, err error) {
	writer.Logf("[STDOUT]: %v", strings.TrimSpace(string(p)))

	return len(p), nil
}

func RedirectOutput(t *testing.T) func() {
	oldOutput := OutputWriter
	OutputWriter = logWriter{t}

	return func() {
		OutputWriter = oldOutput
	}
}

func testRedirectOutput(t *testing.T) {
	assert.Equal(t, os.Stdout, OutputWriter)

	deferFunc := RedirectOutput(t)
	assert.IsType(t, logWriter{}, OutputWriter)
	deferFunc()
	assert.Equal(t, os.Stdout, OutputWriter)
}

// FIXME: check if the json format is valid.
func TestUpdateLanguagesIntregration(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test")
	}

	tmpdir, err := ioutil.TempDir("", "home")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	defer os.RemoveAll(tmpdir)
	defer plusTesting.ChangeEnv("HOME", tmpdir)()
	defer RedirectOutput(t)()

	updateLanguages()
	info, err := os.Stat(path.Join(tmpdir, ".dict_cc", "languages.json"))
	if assert.NoError(t, err) {
		assert.True(t, info.Size() > 100)
	}
}

package httpTesting

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHTTPResponse(t *testing.T) {
	testString := "Hello world :)"
	uri, err := url.Parse("http://google.com")
	if err != nil {
		t.Fatal("Unexpected error parsing the url", err)
	}
	response := GetHTTPResponse(302, uri, testString)

	assert.Equal(t, "302", response.Status)
	assert.Equal(t, 302, response.StatusCode)
	assert.Equal(t, "HTTP/1.0", response.Proto)
	assert.Equal(t, 1, response.ProtoMajor)
	assert.Equal(t, 0, response.ProtoMinor)
	assert.False(t, response.Close)

	content := make([]byte, len(testString))
	response.Body.Read(content)

	assert.Equal(t, []byte(testString), content)

	request := response.Request
	assert.Equal(t, "GET", request.Method)
	assert.Equal(t, "http://google.com", request.URL.String())
}

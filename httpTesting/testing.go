package httpTesting

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// GetHTTPResponse creates a new mocked http response for the given status and content.
func GetHTTPResponse(status int, uri *url.URL, content string) *http.Response {
	contentReader := ioutil.NopCloser(strings.NewReader(content))

	request := &http.Request{
		Method:     "GET",
		URL:        uri,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
	}

	return &http.Response{
		Status:     fmt.Sprint(status),
		StatusCode: status,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Body:       contentReader,
		Close:      false,
		Request:    request,
	}
}

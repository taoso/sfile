package http

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse1(t *testing.T) {
	path := url.PathEscape("/涛叔.html")
	buf := []byte("GET " + path + " HTTP/1.1\r\nA:1\r\nB:2\r\n\r\n")
	req := Request{}
	status, offset := req.Feed(buf)

	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "/涛叔.html", req.Path)
	assert.Equal(t, "1", req.Headers.Get("a"))
	assert.Equal(t, "2", req.Headers.Get("B"))

	assert.Equal(t, len(buf), offset)
	assert.Equal(t, ParseDone, status)
}

func TestParse2(t *testing.T) {
	req := Request{}

	s, o := req.Feed([]byte("GET /foo.html HTTP/1"))
	assert.Equal(t, 14, o)
	assert.Equal(t, ParseVersion, s)

	req.Feed([]byte("HTTP/1.1\r\nA"))
	req.Feed([]byte("A:1\r\nB:2\r\n"))
	s, _ = req.Feed([]byte("\r\n"))

	assert.Equal(t, ParseDone, s)

	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "/foo.html", req.Path)
	assert.Equal(t, "1", req.Headers.Get("a"))
	assert.Equal(t, "2", req.Headers.Get("B"))
}

func TestParse3(t *testing.T) {
	req := Request{}

	req.Feed([]byte("GET /go.mod HTTP/1.1\n"))

	s, _ := req.Feed([]byte("\n"))

	assert.Equal(t, ParseDone, s)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "/go.mod", req.Path)
	assert.Equal(t, "HTTP/1.1", req.Version)
}

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

	req.Feed([]byte("GET /foo.html HTTP/1.1\r\n"))
	req.Feed([]byte("A:1\r\n"))
	req.Feed([]byte("B:2\r\n"))
	s, _ := req.Feed([]byte("\r\n"))

	assert.Equal(t, ParseDone, s)

	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "/foo.html", req.Path)
	assert.Equal(t, "1", req.Headers.Get("a"))
	assert.Equal(t, "2", req.Headers.Get("B"))
}

func TestParse3(t *testing.T) {
	req := Request{}

	req.Feed([]byte("GET /go.mod HTTP/1.1\r\n"))
	req.Feed([]byte("Accept: */*\r\n"))
	req.Feed([]byte("Accept-Encoding: gzip, deflate\r\n"))
	req.Feed([]byte("Connection: keep-alive\r\n"))
	req.Feed([]byte("Host: localhost:8080\r\n"))
	req.Feed([]byte("User-Agent: HTTPie/3.2.1\r\n"))

	s, _ := req.Feed([]byte("\r\n"))

	assert.Equal(t, ParseDone, s)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "/go.mod", req.Path)
	assert.Equal(t, "*/*", req.Headers.Get("Accept"))
	assert.Equal(t, "gzip, deflate", req.Headers.Get("Accept-Encoding"))
	assert.Equal(t, "keep-alive", req.Headers.Get("Connection"))
	assert.Equal(t, "localhost:8080", req.Headers.Get("Host"))
	assert.Equal(t, "HTTPie/3.2.1", req.Headers.Get("User-Agent"))
}
func TestParse4(t *testing.T) {
	req := Request{}

	req.Feed([]byte("GET /go.mod HTTP/1.1\n"))

	s, _ := req.Feed([]byte("\n"))

	assert.Equal(t, ParseDone, s)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "/go.mod", req.Path)
	assert.Equal(t, "HTTP/1.1", req.Version)
}

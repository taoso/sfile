package http

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteChunk(t *testing.T) {
	r := bytes.NewReader([]byte("abc"))
	var w bytes.Buffer
	err := WriteChunk(2, &w, r)

	assert.Nil(t, err)
	assert.Equal(t, "HTTP/1.1 200 OK\r\n"+
		"Transfer-Encoding: chunked\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"\r\n"+
		"2\r\nab\r\n1\r\nc\r\n0\r\n\r\n", w.String())
}

func TestWriteGzip(t *testing.T) {
	r := bytes.NewReader([]byte("aaaaaaaaaa"))
	var w bytes.Buffer
	err := WriteGzip(1, &w, r)

	assert.Nil(t, err)
	assert.Equal(t, "HTTP/1.1 200 OK\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"Content-Length: 27\r\n"+
		"Content-Encoding: gzip\r\n"+
		"\r\n", string(w.Bytes()[:104]))
	zr, err := gzip.NewReader(bytes.NewReader(w.Bytes()[104:]))
	assert.Nil(t, err)
	data, err := io.ReadAll(zr)
	assert.Nil(t, err)

	assert.Equal(t, []byte("aaaaaaaaaa"), data)
}

func TestWriteGzipSmall(t *testing.T) {
	r := bytes.NewReader([]byte("aaaaaaaaaa"))
	var w bytes.Buffer
	err := WriteGzip(20, &w, r)

	assert.Nil(t, err)
	assert.Equal(t, "HTTP/1.1 200 OK\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"Content-Length: 10\r\n"+
		"\r\n"+
		"aaaaaaaaaa", w.String())
}

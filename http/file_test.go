package http

import (
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteChunk(t *testing.T) {
	r := bytes.NewReader([]byte("abc"))
	var w bytes.Buffer
	n, err := WriteChunk(2, &w, r)

	resp := "HTTP/1.1 200 OK\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n" +
		"\r\n" +
		"2\r\nab\r\n1\r\nc\r\n0\r\n\r\n"

	assert.Nil(t, err)
	assert.Equal(t, len(resp), n)
	assert.Equal(t, resp, w.String())
}

func TestWriteGzip(t *testing.T) {
	data := []byte("aaaaaaaaaa")

	var w bytes.Buffer

	zw := gzip.NewWriter(&w)
	zw.Write(data)
	zw.Close()

	resp := append([]byte("HTTP/1.1 200 OK\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"Content-Length: 27\r\n"+
		"Content-Encoding: gzip\r\n"+
		"\r\n"), w.Bytes()...)

	w.Reset()
	r := bytes.NewReader(data)
	n, err := WriteGzip(1, &w, r)

	assert.Nil(t, err)
	assert.Equal(t, len(resp), n)
	assert.Equal(t, resp, w.Bytes())
}

func TestWriteGzipSmall(t *testing.T) {
	r := bytes.NewReader([]byte("aaaaaaaaaa"))
	var w bytes.Buffer
	_, err := WriteGzip(20, &w, r)

	assert.Nil(t, err)
	assert.Equal(t, "HTTP/1.1 200 OK\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"Content-Length: 10\r\n"+
		"\r\n"+
		"aaaaaaaaaa", w.String())
}

package http

import (
	"bytes"
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

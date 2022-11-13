package log

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taoso/sfile/http"
)

func TestLog(t *testing.T) {
	var w bytes.Buffer

	l := &Log{
		EntryNum: 1,
		Interval: 5 * time.Second,
		Writer:   &w,
	}
	l.Init()

	go l.Loop()
	defer l.t.Stop()

	req := http.Request{
		Method:  "GET",
		Path:    "/go.mod",
		Version: "HTTP/1.1",
	}
	req.Headers = http.Headers{}
	req.Headers.Add("User-Agent", "curl/1.1")
	req.Headers.Add("Referer", "https://taoshu.in/")

	l.Log(req, "200", 1024)

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t,
		"GET /go.mod HTTP/1.1 200 1024 https://taoshu.in/ \"curl/1.1\"\n",
		w.String()[len(time.RFC3339)+1:])
}

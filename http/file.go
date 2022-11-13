package http

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func WriteChunk(chunkSize int, w io.Writer, f io.Reader) (int, error) {
	sent := 0
	headerSent := false
	buf := make([]byte, chunkSize)

	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		chunk := buf[:n]
		if !headerSent {
			ctype := http.DetectContentType(chunk)
			i, err := w.Write([]byte("HTTP/1.1 200 OK\r\n" +
				"Transfer-Encoding: chunked\r\n" +
				"Content-Type: " + ctype + "\r\n"))
			if err != nil {
				return 0, err
			}
			sent += i
			headerSent = true
		}
		if n == 0 {
			break
		}
		hexSize := strconv.FormatInt(int64(n), 16)
		i, err := w.Write([]byte("\r\n" + hexSize + "\r\n"))
		if err != nil {
			return 0, err
		}
		sent += i
		if i, err = w.Write(chunk); err != nil {
			return 0, err
		}
		sent += i
	}
	i, err := w.Write([]byte("\r\n0\r\n\r\n"))
	if err != nil {
		return 0, err
	}
	sent += i
	return sent, nil
}

func WriteGzip(zipSize int, w io.Writer, f io.Reader) (int, error) {
	data, err := io.ReadAll(f)
	if err != nil {
		return 0, err
	}

	var gziped bool
	ctype := http.DetectContentType(data)
	if strings.HasPrefix(ctype, "text/") && len(data) > zipSize {
		var buf bytes.Buffer
		zw, _ := gzip.NewWriterLevel(&buf, gzip.DefaultCompression)
		if _, err = zw.Write(data); err != nil {
			return 0, err
		}
		if err = zw.Close(); err != nil {
			return 0, err
		}
		data = buf.Bytes()
		gziped = true
	}

	clen := strconv.Itoa(len(data))
	header := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: " + ctype + "\r\n" +
		"Content-Length: " + clen + "\r\n"
	if gziped {
		header += "Content-Encoding: gzip\r\n\r\n"
	} else {
		header += "\r\n"
	}

	sent, err := w.Write([]byte(header))
	if err != nil {
		return 0, err
	}

	n, err := w.Write(data)
	if err != nil {
		return 0, err
	}
	sent += n
	return sent, nil
}

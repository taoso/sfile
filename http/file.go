package http

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func WriteChunk(chunkSize int, w io.Writer, f io.Reader) error {
	buf := make([]byte, chunkSize)
	headerSent := false
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		chunk := buf[:n]
		if !headerSent {
			ctype := http.DetectContentType(chunk)
			_, err = w.Write([]byte("HTTP/1.1 200 OK\r\n" +
				"Transfer-Encoding: chunked\r\n" +
				"Content-Type: " + ctype + "\r\n"))
			if err != nil {
				return err
			}
			headerSent = true
		}
		if n == 0 {
			break
		}
		hexSize := strconv.FormatInt(int64(n), 16)
		_, err = w.Write([]byte("\r\n" + hexSize + "\r\n"))
		if err != nil {
			return err
		}
		if _, err := w.Write(chunk); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte("\r\n0\r\n\r\n")); err != io.EOF {
		return err
	}
	return nil
}

func WriteGzip(zipSize int, w io.Writer, f io.Reader) error {
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	ctype := http.DetectContentType(data)
	var gziped bool
	if strings.HasPrefix(ctype, "text/") && len(data) > zipSize {
		var buf bytes.Buffer
		zw, _ := gzip.NewWriterLevel(&buf, gzip.DefaultCompression)
		if _, err := zw.Write(data); err != nil {
			return err
		}
		if err := zw.Close(); err != nil {
			return err
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
	_, err = w.Write([]byte(header))
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

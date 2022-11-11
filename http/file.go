package http

import (
	"io"
	"net/http"
	"strconv"
)

func WriteChunk(w io.Writer, f io.Reader) error {
	buf := make([]byte, 1024)
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

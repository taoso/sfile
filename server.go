package sfile

import (
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	shttp "github.com/taoso/sfile/http"
)

type Server struct {
	Root fs.FS
}

func (s *Server) ListenAndServe(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		go s.Serve(c)
	}
}

func (s *Server) Serve(c net.Conn) {
	defer c.Close()

	var n int
	var err error
	var req shttp.Request
	buf := make([]byte, 1024)
	for {
		n, err = c.Read(buf[n:])
		if err != nil {
			return
		}

		status, offset := req.Feed(buf[:n])
		if status == shttp.ParseError {
			return
		} else if status == shttp.ParseDone {
			break
		}
		if offset < n {
			copy(buf, buf[offset:n])
			n -= offset
		} else {
			n = 0
		}
	}

	f, err := s.Root.Open(req.Path[1:])
	if os.IsNotExist(err) {
		c.Write([]byte("HTTP/1.1 404 Not Found\r\nContent-Length:0\r\n\r\n"))
	} else if err != nil {
		msg := err.Error()
		length := strconv.Itoa(len(msg))
		c.Write([]byte("HTTP/1.1 500 Internal Server Error\r\nContent-Length:" + length + "\r\n\r\n" + msg))
	} else {
		buf := make([]byte, 1024)
		headerSent := false
		for {
			n, err := f.Read(buf)
			if err != nil && err != io.EOF {
				log.Println(err)
				return
			}
			chunk := buf[:n]
			if !headerSent {
				ctype := http.DetectContentType(chunk)
				c.Write([]byte("HTTP/1.1 200 OK\r\n" +
					"Transfer-Encoding: chunked\r\n" +
					"Content-Type: " + ctype + "\r\n"))
				headerSent = true
			}
			if n == 0 {
				break
			}
			c.Write([]byte("\r\n" + strconv.FormatInt(int64(n), 16) + "\r\n"))
			if _, err := c.Write(chunk); err != nil {
				log.Println(err)
				return
			}
		}
		c.Write([]byte("\r\n0\r\n\r\n"))
	}
}

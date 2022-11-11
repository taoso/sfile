package sfile

import (
	"io/fs"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/taoso/sfile/file"
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
		go s.serve(c)
	}
}

func (s *Server) serve(c net.Conn) {
	defer c.Close()
	for {
		keepAlive := s.serveOnce(c)
		if !keepAlive {
			return
		}
	}
}

func (s *Server) serveOnce(c net.Conn) bool {
	var n int
	var err error
	var req shttp.Request
	buf := make([]byte, 1024)

	for {
		if n, err = c.Read(buf[n:]); err != nil {
			return false
		}

		status, offset := req.Feed(buf[:n])
		if status == shttp.ParseError {
			return false
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
		c.Write([]byte("HTTP/1.1 404 Not Found\r\n" +
			"Content-Length:0\r\n\r\n"))
	} else if err != nil {
		msg := err.Error()
		length := strconv.Itoa(len(msg))
		c.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n" +
			"Content-Length:" + length + "\r\n\r\n" + msg))
		return false
	} else if err := file.WriteChunk(c, f); err != nil {
		log.Println(err)
		return false
	}

	if req.Headers.Get("Connection") == "close" || req.Version < "HTTP/1.1" {
		return false
	}
	return true
}

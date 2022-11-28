package server

import (
	"io"
	"io/fs"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/taoso/sfile/http"
	slog "github.com/taoso/sfile/log"
)

type Server struct {
	Root        fs.FS
	ChunkSize   int
	GzipSize    int
	ReadTimeout time.Duration

	log *slog.Log
}

func (s *Server) ListenAndServe(addr string) error {
	s.log = &slog.Log{
		EntryNum: 1024,
		Writer:   os.Stderr,
		Interval: 1 * time.Second,
	}
	s.log.Init()
	go s.log.Loop()
	defer s.log.Close()

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
	var r int
	var err error
	var req http.Request
	buf := make([]byte, 60)

	for {
		var n int
		d := time.Now().Add(s.ReadTimeout)
		if err = c.SetReadDeadline(d); err != nil {
			log.Println(err)
			return false
		}

		if n, err = c.Read(buf[r:]); err != nil {
			if !os.IsTimeout(err) && err != io.EOF {
				log.Println(err)
			}
			return false
		}

		status, offset := req.Feed(buf[:n+r])
		if status == http.ParseError {
			log.Println("request parser error")
			return false
		} else if status == http.ParseDone {
			break
		}

		if offset < n {
			copy(buf, buf[offset:n])
			r = n - offset
		} else {
			r = 0
		}
	}

	f, err := s.Root.Open(req.Path[1:])

	if os.IsNotExist(err) {
		sent, _ := c.Write([]byte("HTTP/1.1 404 Not Found\r\n" +
			"Content-Length:0\r\n\r\n"))
		s.log.Log(req, "404", sent)
		return false
	} else if err != nil {
		msg := err.Error()
		length := strconv.Itoa(len(msg))
		sent, _ := c.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n" +
			"Content-Length:" + length + "\r\n\r\n" + msg))
		s.log.Log(req, "500", sent)
		return false
	}

	var sent int
	if strings.Contains(req.Headers.Get("Accept-Encoding"), "gzip") {
		sent, err = http.WriteGzip(s.GzipSize, c, f)
	} else {
		sent, err = http.WriteChunk(s.ChunkSize, c, f)
	}

	if err != nil {
		log.Println(err)
		return false
	}

	s.log.Log(req, "200", sent)

	if req.Headers.Get("Connection") == "close" || req.Version < "HTTP/1.1" {
		return false
	}

	return true
}

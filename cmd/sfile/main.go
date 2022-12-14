package main

import (
	"flag"
	"os"
	"time"

	"github.com/taoso/sfile/server"
)

func main() {
	var root string
	var addr string
	var timeout time.Duration
	var chunk int
	var gzip int

	flag.StringVar(&root, "root", "", "file root path, default is pwd")
	flag.StringVar(&addr, "addr", "127.0.0.1:8080", "listen address and port")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection keep alive timeout")
	flag.IntVar(&chunk, "chunk", 1024, "transfer chunk size")
	flag.IntVar(&gzip, "gzip", 2048, "min file size for gizp")

	flag.Parse()

	var err error
	if root == "" {
		root, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	s := server.Server{}
	s.Root = os.DirFS(root)
	s.ReadTimeout = timeout
	s.ChunkSize = chunk
	s.GzipSize = gzip

	if err := s.ListenAndServe(addr); err != nil {
		panic(err)
	}
}

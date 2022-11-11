package main

import (
	"flag"
	"os"
	"time"

	"github.com/taoso/sfile"
)

func main() {
	var root string
	var addr string
	var timeout time.Duration

	flag.StringVar(&root, "root", "", "file root path, default is pwd")
	flag.StringVar(&addr, "addr", "127.0.0.1:8080", "listen address and port")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection keep alive timeout")

	flag.Parse()

	var err error
	if root == "" {
		root, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	s := sfile.Server{}
	s.Root = os.DirFS(root)
	s.ReadTimeout = timeout

	if err := s.ListenAndServe(addr); err != nil {
		panic(err)
	}
}

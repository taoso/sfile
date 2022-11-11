package main

import (
	"flag"
	"os"

	"github.com/taoso/sfile"
)

func main() {
	var root string
	flag.StringVar(&root, "root", "", "file root path")
	flag.Parse()

	if root == "" {
		root, _ = os.Getwd()
	}

	s := sfile.Server{}
	s.Root = os.DirFS(root)
	err := s.ListenAndServe("127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
}

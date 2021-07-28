package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"dsqlite/server"
)

// variable
var host string
var port int
var path string
var header string

func init() {
	flag.StringVar(&host, "localhost", "127.0.0.1", "dsqlite's server address")
	flag.IntVar(&port, "port", 4001, "dsqlite's server port")
	flag.StringVar(&header, "join", "", "dsqlite's header address")
	// usage for flag
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: [arguments] <data-path>\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	path = flag.Arg(0)
	// create workspace dir.
	if err := os.MkdirAll(path, 0774); err != nil {
		log.Fatalf("Unable to create path: %v", err)
		return
	}

	// create dsqlite server
	s := server.NewServer(host, port, path)
	log.Fatal(s.ListenAndLeave(header))
}

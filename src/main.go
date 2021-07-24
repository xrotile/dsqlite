package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"server"
)

// variable
var host string
var port int
var path string

func init() {
	flag.StringVar(&host, "localhost", "127.0.0.1", "dsqlite's server address")
	flag.IntVar(&port, "port", 4001, "dsqlite's server port")
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
	s := server.NewServer(host, port)
	log.Fatal(s.ListenAndLeave())
}

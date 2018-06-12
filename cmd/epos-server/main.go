package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/cloudinn/escpos"
)

var (
	flagListen   = flag.String("l", "127.0.22.8:80", "listen")
	flagEndpoint = flag.String("endpoint", escpos.DefaultEndpoint, "endpoint")
	flagPrinter  = flag.String("p", "", "path to printer")
)

func main() {
	flag.Parse()

	if *flagPrinter == "" {
		log.Fatal("must specify path to printer via -p")
	}

	// open printer
	f, err := os.Create(*flagPrinter)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// create printer
	ep, err := escpos.NewPrinter(f)
	if err != nil {
		log.Fatal(err)
	}

	// create server
	s, err := escpos.NewServer(ep)
	if err != nil {
		log.Fatal(err)
	}

	// set up mux
	mux := http.NewServeMux()
	mux.Handle(*flagEndpoint, s)
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	log.Fatal(http.ListenAndServe(*flagListen, mux))
}

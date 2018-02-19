package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/knq/escpos"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
)

var listenAddr = flag.String("l", "127.0.22.8", "Address to listen on")
var port = flag.Int("port", 80, "Port to listen on")
var printerPath = flag.String("p", "/dev/usb/lp0", "Path to printer")

type EposServer struct {
	r             *mux.Router
	printer       *escpos.Escpos
	printerWriter *bufio.Writer
}

func writeSoapResponse(rw http.ResponseWriter, req *http.Request, code string) {
	success_str := "false"
	if code == "" {
		success_str = "true"
	}

	response := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
  <s:Body xmlns:m="http://www.epson-pos.com/schemas/2011/03/epos-print">
	<m:response success="%s" code="%s" status="0"></m:response>
  </s:Body>
</s:Envelope>`, success_str, code)

	log.Printf("Sending:\n%s\n", response)

	// inject response
	rw.Header().Set("Content-Type", req.Header.Get("Content-Type"))
	rw.Write([]byte(response))
}

func getEposNodes(doc *xml.XmlDocument) (retnodes []xml.Node, err error) {
	// grab the 'Body' element
	path := xpath.Compile("*[local-name()='Body']")
	nodes, e := doc.Root().Search(path)
	if e != nil {
		err = e
		return
	}

	// check that the data is present
	if len(nodes) < 1 || nodes[0].CountChildren() < 1 {
		err = errors.New("bad data")
		return
	}

	// get epos data
	return nodes[0].FirstChild().Search("./*")
}

func (s *EposServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// send origin headers
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, If-Modified-Since, SOAPAction")
	}

	// stop if its options
	if req.Method == "OPTIONS" {
		log.Printf("OPTIONS %s\n", req.URL)
		return
	}

	// handle crappy soap action
	if req.Method == "POST" {
		// grab posted body
		data, _ := ioutil.ReadAll(req.Body)
		log.Printf("POST %s:\n%s\n\n", req.URL, string(data))

		// parse xml with gokogiri
		doc, _ := gokogiri.ParseXml(data)
		defer doc.Free()

		// load print nodes from xml doc
		epos_nodes, err := getEposNodes(doc)
		if err != nil {
			rw.WriteHeader(503)
			log.Fatal(err)
			return
		}

		// init printer
		s.printer.Init()

		// loop over nodes
		for _, en := range epos_nodes {
			// grab name and inner text
			name := en.Name()
			content := en.Content()

			// grab parameters
			params := make(map[string]string)
			for _, attr := range en.Attributes() {
				params[attr.Name()] = attr.Value()
			}

			// write data to printer
			s.printer.WriteNode(name, params, content)
		}

		// end
		s.printer.End()

		// flush writer
		s.printerWriter.Flush()

		//rw.WriteHeader(402)
		// write soap response
		writeSoapResponse(rw, req, "")

		return
	}

	// force an error for everything else
	rw.WriteHeader(403)

	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

func main() {
	flag.Parse()

	// open printer
	f, err := os.Create(*printerPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// setup buffered writer
	w := bufio.NewWriter(f)
	ep := escpos.New(w)

	// set up service router
	r := mux.NewRouter()
	http.Handle("/cgi-bin/epos/service.cgi", &EposServer{r, ep, w})

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *listenAddr, *port), nil))
}

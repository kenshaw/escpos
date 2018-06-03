package connection

import (
	"cloudinn/escpos"
	"io"
	"log"
	"net"
	"os"
)

var f io.Writer
var err error

//NewConnection creats a connection with a usb printer or a network printer and
//returns an object to use escops package functions with
func NewConnection(connectionType string, connectionHost string) *escpos.Printer {

	if connectionType == "usb" {
		f, err = os.OpenFile(connectionHost, os.O_WRONLY, 0)
	} else if connectionType == "network" {
		f, err = net.Dial("tcp", connectionHost)

	}
	if err != nil {
		log.Fatal(err)
	}
	printerObj, err := escpos.NewPrinter(f)
	if err != nil {
		log.Fatal(err)
	}
	return printerObj

}

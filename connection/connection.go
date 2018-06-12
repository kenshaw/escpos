package connection

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/cloudinn/escpos"
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
		log.Println(err.Error())
	}
	printerObj, err := escpos.NewPrinter(f)
	if err != nil {
		log.Println(err.Error())
	}
	return printerObj

}

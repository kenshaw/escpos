package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"

	_ "image/jpeg"

	"github.com/kenshaw/escpos"
	"github.com/kenshaw/escpos/raster"
)

var (
	lpDev    = flag.String("p", "/dev/usb/lp0", "Printer dev file")
	maxWidth = flag.Int("printer-max-width", 512, "Printer max width in pixels")

	ep *escpos.Printer
)

func main() {
	flag.Parse()

	f, err := os.OpenFile(*lpDev, os.O_WRONLY, 0)
	if err != nil {
		log.Fatal(err)
	}

	ep, err = escpos.NewPrinter(f)
	if err != nil {
		log.Fatal(err)
	}

	ep.Init()

	defer func() {
		ep.Cut()
		ep.End()
	}()

	ep.Text(nil, "sample text...\n\n")

	for _, font := range []string{"A", "B", "C"} {
		ep.SetFont(font)
		ep.Text(nil, fmt.Sprintf("sample text, font %s...\n\n", font))
	}
	ep.SetFont("B")

	for _, format := range []int{0, 1, 2, 3, 4, 73} {
		ep.Text(nil, fmt.Sprintf("sample barcode, format %d:\n", format))
		ep.Barcode("123456", format)
		ep.Linefeed()
	}

	ep.Text(nil, "sample image:\n")
	rasterImage()

	ep.Text(nil, "cash code:\n\n")
	ep.Cash()
}

func rasterImage() {
	imgFile, err := os.Open("crab_nebula.jpg")
	if err != nil {
		log.Fatal(err)
	}

	img, imgFormat, err := image.Decode(imgFile)
	imgFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Loaded image, format: ", imgFormat)

	rasterConv := &raster.Converter{
		MaxWidth:  *maxWidth,
		Threshold: 0.5,
	}

	rasterConv.Print(img, ep)
}

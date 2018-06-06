package main

import (
	"flag"
	"image"
	"log"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/cloudinn/escpos"
	"github.com/cloudinn/escpos/raster"
)

var (
	lpDev     = flag.String("p", "/dev/usb/lp0", "Printer dev file")
	imgPath   = flag.String("i", "image.png", "Input image")
	threshold = flag.Float64("t", 0.5, "Black/white threshold")
	align     = flag.String("a", "center", "Alignment (left, center, right)")
	doCut     = flag.Bool("c", false, "Cut after print")
	maxWidth  = flag.Int("printer-max-width", 512, "Printer max width in pixels")
)

func main() {
	flag.Parse()

	imgFile, err := os.Open(*imgPath)
	if err != nil {
		log.Fatal(err)
	}

	img, imgFormat, err := image.Decode(imgFile)
	imgFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Loaded image, format: ", imgFormat)

	// ----------------------------------------------------------------------

	f, err := os.OpenFile(*lpDev, os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	log.Print(*lpDev, " open.")

	ep := escpos.New(f)

	ep.Init()

	ep.SetAlign(*align)

	rasterConv := &raster.Converter{
		MaxWidth:  *maxWidth,
		Threshold: *threshold,
	}

	rasterConv.Print(img, ep)

	if *doCut {
		ep.Cut()
	}
	ep.End()
}

package raster

import (
	"image"
	"image/color"
)

type Target interface {
	Raster(width, height, bytesWidth int, rasterData []byte)
}

type Converter struct {
	// The maximum line width of the printer, in dots
	MaxWidth int

	// The threashold between white and black dots
	Threshold float64
}

func (c *Converter) Print(img image.Image, target Target) {
	sz := img.Bounds().Size()

	data, rw, bw := c.ToRaster(img)

	target.Raster(rw, sz.Y, bw, data)
}

func (c *Converter) ToRaster(img image.Image) (data []byte, imageWidth, bytesWidth int) {
	sz := img.Bounds().Size()

	// lines are packed in bits
	imageWidth = sz.X
	if imageWidth > c.MaxWidth {
		// truncate if image is too large
		imageWidth = c.MaxWidth
	}

	bytesWidth = imageWidth / 8
	if imageWidth%8 != 0 {
		bytesWidth += 1
	}

	data = make([]byte, bytesWidth*sz.Y)

	for y := 0; y < sz.Y; y++ {
		for x := 0; x < imageWidth; x++ {
			if lightness(img.At(x, y)) >= c.Threshold {
				// position in data is: line_start + x / 8
				// line_start is y * bytesWidth
				// then 8 bits per byte
				data[y*bytesWidth+x/8] |= 0x80 >> uint(x%8)
			}
		}
	}

	return
}

const (
	lumR, lumG, lumB = 55, 182, 18
)

func lightness(c color.Color) float64 {
	r, g, b, _ := c.RGBA()

	return float64(lumR*r+lumG*g+lumB*b) / float64(0xffff*(lumR+lumG+lumB))
}

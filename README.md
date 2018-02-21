# About escpos #

This is a simple [Go][1] package that provides [ESC-POS][2] library functions
to help with sending control codes to a ESC-POS capable printer such as an
Epson TM-T82 or similar.

These printers are often used in retail environments in conjunction with a
point-of-sale (POS) system.

## Installation ##

Install the package via the following:

    go get -u github.com/kenshaw/escpos

## Example epos-server ##

An example EPOS server implementation is available in the [cmd/epos-server][3]
subdirectory of this project. This example server is more or less compatible
with [Epson TM-Intelligent][4] printers and print server implementations.

## Usage ##

The escpos package can be used similarly to the following:

```go
package main

import (
    "bufio"
    "os"

    "github.com/kenshaw/escpos"
)

func main() {
    f, err := os.Open("/dev/usb/lp3")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    w := bufio.NewWriter(f)
    p := escpos.New(w)

    p.Init()
    p.SetSmooth(1)
    p.SetFontSize(2, 3)
    p.SetFont("A")
    p.Write("test ")
    p.SetFont("B")
    p.Write("test2 ")
    p.SetFont("C")
    p.Write("test3 ")
    p.Formfeed()

    p.SetFont("B")
    p.SetFontSize(1, 1)

    p.SetEmphasize(1)
    p.Write("halle")
    p.Formfeed()

    p.SetUnderline(1)
    p.SetFontSize(4, 4)
    p.Write("halle")

    p.SetReverse(1)
    p.SetFontSize(2, 4)
    p.Write("halle")
    p.Formfeed()

    p.SetFont("C")
    p.SetFontSize(8, 8)
    p.Write("halle")
    p.FormfeedN(5)

    p.Cut()
    p.End()

    w.Flush()
}
```

## TODO
- Fix barcode/image support

[1]: http://www.golang.org/project
[2]: https://en.wikipedia.org/wiki/ESC/P
[3]: cmd/epos-server
[4]: https://c4b.epson-biz.com

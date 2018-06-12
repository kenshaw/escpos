# About epos-server #

This is a quick and dirty Golang implementation of a [Epson TM-Intelligent][1]
print server. This also serves as example code for the [escpos][2] package.
This is "more-or-less" compatible with the ePOS-Print API and shows how
ePOS-XML is translated into simple ESCPOS data.

This has been tested and works as expected on Linux.

## Installation and Building ##

You will likely need the libxml2 and related headers to build the gokogiri
dependency.

If they are not already installed on your system, you will need to do something
like the following for your system:

    sudo aptitude install libxml2 libxml2-dev

Then install via the following:

    go get -u github.com/kenshaw/escpos/epos-server

You should then be able to build the epos-server like this:

    go build github.com/kenshaw/escpos/epos-server

## Usage ##

You can specify the address and port to listen on, as well as the path to the
printer:

    user@host# ./epos-server --help
    Usage of ./epos-server:
      -l string
            Address to listen on (default "127.0.22.8")
      -p string
            Path to printer (default "/dev/usb/lp0")
      -port int
            Port to listen on (default 80)

## TODO ##

The following still needs to be implemented:

* Fix image decoding and printing

[1]: https://c4b.epson-biz.com/
[2]: https://github.com/kenshaw/escpos

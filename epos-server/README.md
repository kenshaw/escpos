# About epos-server #

This is a quick and dirty Golang implementation of a [Epson TM-Intelligent](https://c4b.epson-biz.com/) 
print server. This also serves as example code for the
[escpos](https://github.com/knq/escpos) package.

This has been tested and works as expected on Linux.

## Installation and Building ##

You will likely need the libxml2 and related headers to build the gokogiri
dependency.

If they are not already installed on your system, you will need to do something
like the following for your system:

    sudo aptitude install libxml2 libxml2-dev

You should then be able to install via the following:

    go get -u github.com/knq/escpos/epos-server


You should then be able to build the epos-server like the following:

    go build github.com/knq/escpos/epos-server

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

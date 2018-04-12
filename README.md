# wt521-rs232-web-interface
Web interface for the Vaisala WT521 using RS-232.

## Installation
```sh
go get github.com/asmundstavdahl/wt521-rs232-web-interface
go install github.com/asmundstavdahl/wt521-rs232-web-interface
```

## Usage
```sh
$ wt521-rs232-web-interface -help
Usage of wt521-rs232-web-interface:
  -baud int
    	baud rate (WT521's facotry default is 1200) (default 1200)
  -device string
    	serial port to use (default "/dev/ttyS0")
  -port int
    	port to open for HTTP server (default 8080)

```

## Credits
Cheers to Wikimedia user [El Grafo](https://commons.wikimedia.org/wiki/User:El_Grafo) for making the [compass SVG](https://en.wiktionary.org/wiki/File:Compass-icon_bb_NEbE.svg) and for dedicating it to the public domain.

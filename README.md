# wt521-rs232-web-interface ![ci_build](https://travis-ci.org/asmundstavdahl/wt521-rs232-web-interface.svg?branch=master)
Web interface for the Vaisala WT521 using RS-232.

## Installation
```sh
# Download
git clone https://github.com/asmundstavdahl/wt521-rs232-web-interface.git
cd wt521-rs232-web-interface/

# Compile
go build
# Run
./wt521-rs232-web-interface
## OR
# Compile and run
make
```

## Usage
```sh
./wt521-rs232-web-interface -help
```
### Proxy configuration
```sh
# Public side:
./wt521-rs232-web-interface -proxy -port=8081

# WT521 side:
./wt521-rs232-web-interface -device=/dev/ttyS0 -report-to=http://localhost:8081

# Now the WT521 side will send weather data over HTTP to the public side, so that
# 1. the public side does not need to be directly connected to the WT521, and
# 2. the WT521 side can be behind a firewall with only outbound network access
```

## Credits
Cheers to Wikimedia user [El Grafo](https://commons.wikimedia.org/wiki/User:El_Grafo) for making the [compass SVG](https://en.wiktionary.org/wiki/File:Compass-icon_bb_NEbE.svg) and for dedicating it to the public domain.

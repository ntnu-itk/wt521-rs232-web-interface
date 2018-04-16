package main

import "flag"

var flagVerbose bool

func init() {
	flag.BoolVar(&flagVerbose, "verbose", false, "turn on verbose logging")
}

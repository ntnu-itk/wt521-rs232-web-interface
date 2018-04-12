package main

import (
	"flag"
	"log"

	"github.com/tarm/serial"
)

var flagDevice string
var flagBaud int

func init() {
	flag.StringVar(&flagDevice, "device", "/dev/ttyS0", "serial port to use")
	flag.IntVar(&flagBaud, "baud", 1200, "baud rate (WT521's factory default is 1200)")
}

func openSerialPort() *serial.Port {
	serialPort, err := serial.OpenPort(
		&serial.Config{
			Name:     flagDevice,
			Baud:     flagBaud,
			Parity:   serial.ParityNone,
			StopBits: serial.Stop1,
			Size:     8})
	if err != nil {
		log.Fatal(err)
	}

	return serialPort
}

func serialMonitor(serialPort *serial.Port, patchChannel PatchChannel) {
	patch := &StatePatch{}

	for true {
		patch.parse(serialPort)
		select {
		case <-patchChannel:
			log.Println("Discarded a state patch. The stateKeeper goroutine may have failed.")
		default:
			//log.Println("No patch in queue")
		}
		patchChannel <- *patch
	}
}

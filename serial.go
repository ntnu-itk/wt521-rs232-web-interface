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
	flag.IntVar(&flagBaud, "baud", 9600, "baud rate (WT521's facotry default is 1200 but it should be reconfigured to 9600)")
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
			log.Println("Discarded a patch. The stateKeeper goroutine may have crashed.")
		default:
			//log.Println("No patch in queue")
		}
		patchChannel <- *patch
	}
}

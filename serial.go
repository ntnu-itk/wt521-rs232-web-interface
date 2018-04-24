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
		log.Fatalf("[openSerialPort] (FATAL) %s", err)
	}

	return serialPort
}

func SerialReader(serialPort *serial.Port, byteChannel chan<- byte) {
	buf := make([]byte, 1)
	for {
		n, err := serialPort.Read(buf)
		if err != nil {
			log.Printf("[SerialReader] Could not read byte from serial port; %v", err)
		} else {
			if n == 0 {
				log.Printf("[SerialReader] Zero bytes read from serial port but no error.")
			} else {
				byteChannel <- buf[0]
			}
		}
	}
}

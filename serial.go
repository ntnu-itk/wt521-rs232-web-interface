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

func SerialMonitor(serialPort *serial.Port, byteSubscribers []chan byte) {
	buf := make([]byte, 1)
	for {
		n, err := serialPort.Read(buf)
		if err != nil {
			log.Printf("Could not read byte from serial port; %v", err)
		} else {
			if n == 0 {
				log.Printf("Zero bytes read from serial port but no error.")
			} else {
				for i := 0; i < len(byteSubscribers); i++ {
					if flagVerbose {
						log.Printf("Sent byte 0x%02X to subscriber %d", buf[0], 1+i)
					}
					byteSubscribers[i] <- buf[0]
				}
			}
		}
	}
}

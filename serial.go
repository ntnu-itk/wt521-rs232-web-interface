package main

import (
	"flag"
	"fmt"
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

func SerialMonitor(serialPort *serial.Port, patchChannel PatchChannel) {
	patch := &StatePatch{}

	for true {
		patch.Parse(serialPort)
		select {
		case <-patchChannel:
			log.Println("Discarded a state patch. The stateKeeper goroutine may have failed.")
		default:
			//log.Println("No patch in queue")
		}
		patchChannel <- *patch
	}
}

func (patch *StatePatch) Parse(serialPort *serial.Port) error {
	b := make([]byte, 500)
	for b[0] != byte('$') {
		serialPort.Read(b[0:1])
	}

	bytesRead := 0
	for {
		n, err := serialPort.Read(b[bytesRead:])
		bytesRead += n
		if err != nil {
			log.Printf("Error while reading data bytes: %s", err)
		}
		if rune(b[bytesRead-1]) == '*' {
			//log.Printf("Have read %d bytes: %s", bytesRead, b)
			break
		}
		if bytesRead == cap(b) {
			log.Printf("Warning: read the maximum of %d bytes into buffer but no delimiter found; bailing early. If this happens consistently, check if serial port is configured correctly.", bytesRead)
			log.Printf("         our buffer contains this: %s", string(b))
			break
		}
	}

	if flagVerbose {
		log.Println(string(b))
	}

	n, err := fmt.Sscanf(string(b),
		"WIMWV,%d,R,%f,M,A*",
		&patch.windAngle,
		&patch.windSpeed)

	if n != 2 || err != nil {
		log.Printf("Failed to extract data from buffer. Was able to parse %d of 2 numbers of interest with error '%s'", n, err)
	}

	return err
}

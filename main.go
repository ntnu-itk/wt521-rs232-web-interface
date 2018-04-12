package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tarm/serial"
)

type State struct {
	windAngle   int16
	windSpeed   float64
	lastUpdated time.Time
}

type StatePatch State

type PatchChannel chan StatePatch
type StateChannel chan State
type RequestChannel chan bool

var patchChannel PatchChannel
var stateChannel StateChannel
var requestChannel RequestChannel

var flagName string
var flagBaud int
var flagPort int

func openSerialPort() (*serial.Port, error) {
	return serial.OpenPort(
		&serial.Config{
			Name:     flagName,
			Baud:     flagBaud,
			Parity:   serial.ParityNone,
			StopBits: serial.Stop1,
			Size:     8})
}

func handleFlags() {
	flag.StringVar(&flagName, "dev", "/dev/ttyS0", "serial port to use")
	flag.IntVar(&flagBaud, "baud", 9600, "baud rate (WT521's facotry default is 1200 but it should be reconfigured to 9600)")
	flag.IntVar(&flagPort, "port", 8080, "port to open for HTTP server")
}

func main() {
	patchChannel = make(PatchChannel, 1)
	stateChannel = make(StateChannel, 0)
	requestChannel = make(RequestChannel, 0)

	handleFlags()

	serialPort, err := openSerialPort()
	if err != nil {
		log.Fatal(err)
	}

	go stateKeeper(patchChannel, requestChannel, stateChannel)
	go serialMonitor(serialPort, patchChannel)

	log.Fatal(httpServer(requestChannel, stateChannel))
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

func stateKeeper(patchChannel PatchChannel, requestChannel RequestChannel, stateChannel StateChannel) {
	state := State{}
	patch := StatePatch{}

	for {
		select {
		case patch = <-patchChannel:
			state.patch(patch)
		case <-requestChannel:
			stateChannel <- state
		}
		log.Println(&state)
	}
}

func httpServer(requestChannel RequestChannel, stateChannel StateChannel) error {
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		requestChannel <- true
		state := <-stateChannel
		log.Printf("Got state %s\n", state.String())

		jsonString := fmt.Sprintf(`{
    "speed": %.1f,
    "angle": %d,
    "updated": "%s"
}`,
			state.windSpeed,
			state.windAngle,
			state.lastUpdated)

		log.Println(jsonString)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, jsonString)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestChannel <- true
		state := <-stateChannel
		log.Printf("Got state %s\n", state.String())

		jsonString := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8" />
</head>
<body>
    <h1>Værinfo</h1>
    Speed: %.1f
    <br>
    Angle: %d
    <br>
    Updated: "%s
    <script>setTimeout('location.reload()', 1010)</script>`,
			state.windSpeed,
			state.windAngle,
			state.lastUpdated)

		log.Println(jsonString)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, jsonString)
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", flagPort), nil)
}

func (patch *StatePatch) parse(serialPort *serial.Port) error {
	b := make([]byte, 500)
	for b[0] != byte('$') {
		serialPort.Read(b[0:1])
	}

	bytesRead := 0
	for {
		n, err := serialPort.Read(b[bytesRead:])
		bytesRead += n
		if err != nil {
			log.Printf("Error reading 'till newline: %s", err)
		}
		if rune(b[bytesRead-1]) == '*' {
			log.Printf("Have read %d bytes: %s", bytesRead, b)
			break
		}
		if bytesRead == cap(b) {
			log.Printf("Warning: read the maximum of %d bytes into buffer but no delimiter found; bailing early. If this happens consistently, check if serial port is configured correctly.", bytesRead)
			log.Printf("Warning (cont.): our buffer contains this: %s", string(b))
			break
		}
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

func (state *State) patch(patch StatePatch) {
	state.windSpeed = patch.windSpeed
	state.windAngle = patch.windAngle
	state.lastUpdated = time.Now()
}

func (state *State) String() string {
	return fmt.Sprintf("State{ Speed:%f  Angle:%d  Updated:%s }",
		state.windSpeed,
		state.windAngle,
		state.lastUpdated)
}

func (patch *StatePatch) String() string {
	return fmt.Sprintf("StatePatch{ Speed:%f  Angle:%d }",
		patch.windSpeed,
		patch.windAngle)
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/tarm/serial"
)

var flagPort int

func init() {
	flag.IntVar(&flagPort, "port", 8080, "port to open for HTTP server")
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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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

		jsonString := fmt.Sprintf(`{
    "speed": %.1f,
    "angle": %d,
    "updated": "%s"
}`,
			state.windSpeed,
			state.windAngle,
			SimpleTimeString(state.lastUpdated))

		log.Printf("%s => serve JSON of %s", r.URL.Path, state.String())
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, jsonString)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fileToServe := "www/" + r.URL.Path
		fileToServe = strings.Replace(fileToServe, "..", ".", -1)
		fileToServe = strings.Replace(fileToServe, "//", "/", -1)

		if r.URL.Path == "/" || r.URL.Path == "index.html" {
			fileToServe = "www/home.html"
		}

		content, err := ioutil.ReadFile(fileToServe)
		if err != nil {
			log.Printf("Error serving static file %s: %s", fileToServe, err)
		}

		log.Printf("%s => serve static file %s (length %d)", r.URL.Path, fileToServe, len(content))

		filenameDotParts := strings.Split(fileToServe, ".")
		if len(filenameDotParts) > 1 {
			mimeType := "text/html"

			switch filenameDotParts[1] {
			case "js":
				mimeType = "text/javascript"
			case "css":
				mimeType = "text/css"
			case "svg":
				mimeType = "image/svg+xml"
			case "png":
				mimeType = "image/png"
			}

			w.Header().Set("Content-Type", mimeType)
		}
		w.Write(content)
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

	n, err := fmt.Sscanf(string(b),
		"WIMWV,%d,R,%f,M,A*",
		&patch.windAngle,
		&patch.windSpeed)

	if n != 2 || err != nil {
		log.Printf("Failed to extract data from buffer. Was able to parse %d of 2 numbers of interest with error '%s'", n, err)
	}

	return err
}

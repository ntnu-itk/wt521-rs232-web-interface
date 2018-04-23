package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	homeHtmlFilePath = "www/home.html"
)

var flagPort int
var flagPollTimeout int64

func init() {
	flag.IntVar(&flagPort, "port", 8080, "port to open for HTTP server")
	flag.Int64Var(&flagPollTimeout, "poll-timeout", 100, "seconds to wait for new state before defaulting to the previous one when client is long polling")
}

func HttpServer(
	stateRequestChannel chan *StateRequest,
	currentStateChannel chan State,
	stateHistory *StateHistory) error {
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		httpHandleJSON(w, r, stateRequestChannel, currentStateChannel)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpHandleRoot(w, r, stateHistory)
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", flagPort), nil)
}

func httpHandleRoot(w http.ResponseWriter,
	r *http.Request,
	stateHistory *StateHistory) {
	fileToServe := "www/" + r.URL.Path
	fileToServe = strings.Replace(fileToServe, "..", ".", -1)
	fileToServe = strings.Replace(fileToServe, "//", "/", -1)

	if r.URL.Path == "/" || r.URL.Path == "index.html" {
		fileToServe = homeHtmlFilePath
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

	if fileToServe == homeHtmlFilePath {
		w.Write([]byte(processHomeFile(string(content), stateHistory)))
	} else {
		w.Write(content)
	}
}

func httpHandleJSON(w http.ResponseWriter,
	r *http.Request,
	stateRequestChannel chan<- *StateRequest,
	stateChannel <-chan State) {

	var state State

	_, longPollMode := r.URL.Query()["wait"]
	if longPollMode {
		select {
		case state = <-stateChannel:
			log.Printf("Waited for state %s", state.String())
			w.WriteHeader(http.StatusOK)
			break
		case <-time.After(time.Duration(flagPollTimeout) * time.Second):
			log.Printf("Long poll timed out after %d seconds", flagPollTimeout)
			w.WriteHeader(http.StatusNotModified)
			break
		}
	} else {
		stateRequest := NewStateRequest()
		stateRequestChannel <- stateRequest
		state = <-stateRequest.reply
	}

	jsonString := fmt.Sprintf(`{
    "speed": %.1f,
    "angle": %d,
    "time": "%s"
}`,
		state.windSpeed,
		state.windAngle,
		SimpleTimeString(state.lastUpdated))

	log.Printf("%s => serve JSON of %s", r.URL.Path, state.String())
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonString))
}

func processHomeFile(content string, stateHistory *StateHistory) string {
	return strings.Replace(
		content,
		"__SAMPLES__",
		stateHistory.ToJSON(),
		-1)
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"mime"
	"net/http"
	"strings"
	"time"
)

const (
	indexHtmlFilePath = "www/index.html"
)

var flagPort int
var flagCameraUrl string
var flagCameraHref string
var flagInterval float64
var flagDaysOfHistory float64

// Max number of samples to keep before discarding oldest
var flagHistorySamples int
var flagPollTimeout int64

func init() {
	flag.IntVar(&flagPort, "port", 8080,
		"port to open for HTTP server")
	flag.StringVar(&flagCameraUrl, "camera-url", "https://www.vegvesen.no/public/webkamera/kamera?id=110409",
		"URL to use for the src attribute of the image in the top left corner")
	flag.StringVar(&flagCameraHref, "camera-href", "",
		"HREF of the link of the image in the top left corner; empty value => use camera-url")
	flag.Float64Var(&flagInterval, "interval", 3,
		"should be the same as the MWV interval of the WT521, see setup.md")
	flag.Float64Var(&flagDaysOfHistory, "history-days", 1,
		"days worth of sample data to keep and show")

	flagHistorySamples = int(math.Round((flagDaysOfHistory * 24 * 3600) / 3))
	flagPollTimeout = int64(math.Round(flagInterval * 1000 * 1.1))
}

type tplData struct {
	StateHistory *StateHistory
	WebcamURL    string
	WebcamHref   string
	MaxLines     int
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

	if r.URL.Path == "/" {
		fileToServe = indexHtmlFilePath
	}

	content, err := ioutil.ReadFile(fileToServe)
	if err != nil {
		log.Printf("[HttpServer] Error serving static file %s: %s", fileToServe, err)
	}

	if flagVerbose {
		log.Printf("[HttpServer] %s => serve static file %s (length %d)", r.URL.Path, fileToServe, len(content))
	}

	mimeType := setMimeType(w, fileToServe)

	if flagVerbose {
		log.Printf("[HttpServer] MIME type is %s", mimeType)
	}

	if strings.Index(mimeType, "text/html") >= 0 {
		var tplCameraHref string
		if flagCameraHref == "" {
			tplCameraHref = flagCameraUrl
		} else {
			tplCameraHref = flagCameraHref
		}

		data := tplData{
			StateHistory: stateHistory,
			WebcamURL:    flagCameraUrl,
			WebcamHref:   tplCameraHref,
			MaxLines:     flagHistorySamples}

		t, err := template.ParseGlob("www/*.html")

		if err != nil {
			log.Printf("[HttpServer] Error at ParseGlob: %s", err)
		} else {
			err = t.ExecuteTemplate(w, "index.html", data)

			if err != nil {
				log.Printf("[HttpServer] Error executing template: %s", err)
			}
		}
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
			if flagVerbose {
				log.Printf("[HttpServer] Waited for state %s", state.String())
			}
			w.WriteHeader(http.StatusOK)
			break
		case <-time.After(time.Duration(flagPollTimeout) * time.Millisecond):
			log.Printf("[HttpServer] Long poll timed out after %d milliseconds, sending previous state", flagPollTimeout)
			w.WriteHeader(http.StatusNotModified)
			break
		}
	} else {
		stateRequest := NewStateRequest()
		stateRequestChannel <- stateRequest
		state = <-stateRequest.reply
	}

	var jsonString string
	jsonBytes, err := json.Marshal(state)
	if err != nil {
		log.Printf("[HttpServer] Could not marshal state: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	} else {
		jsonString = string(jsonBytes)
	}

	if flagVerbose {
		log.Printf("[HttpServer] %s => serve JSON of %s", r.URL.Path, state.String())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonString))
}

func setMimeType(w http.ResponseWriter, fileToServe string) (mimeType string) {
	filenameDotParts := strings.Split(fileToServe, ".")
	mimeType = "text/plain"
	if len(filenameDotParts) > 1 {
		mimeType = mime.TypeByExtension(
			fmt.Sprintf(
				".%s",
				filenameDotParts[len(filenameDotParts)-1]))
	}

	w.Header().Set("Content-Type", mimeType)

	return
}

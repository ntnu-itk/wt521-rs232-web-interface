package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"strings"
	"time"
)

const (
	indexHtmlFilePath = "www/index.html"
)

var flagPort int
var flagPollTimeout int64
var flagCameraUrl string

func init() {
	flag.IntVar(&flagPort, "port", 8080,
		"port to open for HTTP server")
	flag.Int64Var(&flagPollTimeout, "poll-timeout", 100,
		"seconds to wait for new state before defaulting to the previous one when client is long polling")
	flag.StringVar(&flagCameraUrl, "camera-url", "https://www.vegvesen.no/public/webkamera/kamera?id=110409",
		"URL to use for the src attribute of the image in the top left corner")
}

type tplData struct {
	StateHistory *StateHistory
	WebcamURL    string
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
		log.Printf("MIME type is %s", mimeType)
	}

	if strings.Index(mimeType, "text/html") >= 0 {
		data := tplData{
			StateHistory: stateHistory,
			WebcamURL:    flagCameraUrl}

		t, err := template.ParseGlob("www/*.html")

		if err != nil {
			log.Printf("[HttpServer] Error at ParseGlob: %s", err)
		} else {
			err = t.ExecuteTemplate(w, "index.html", data)

			if err != nil {
				log.Printf("Error executing template: %s", err)
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
		case <-time.After(time.Duration(flagPollTimeout) * time.Second):
			log.Printf("[HttpServer] Long poll timed out after %d seconds, sending previous state", flagPollTimeout)
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

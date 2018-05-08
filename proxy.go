package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"
)

var flagEnableProxy bool

func init() {
	flag.BoolVar(&flagEnableProxy, "proxy", false, "listen for state reports from an instance of this program run with -report-to=<URL>")
}

func ConfigureProxy(patchChannel chan<- StatePatch) {
	handler := &proxyHandler{patchChannel: patchChannel}
	http.Handle("/report/json", handler)
}

type proxyHandler struct {
	patchChannel chan<- StatePatch
}

func (ph *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("[Proxy] Error parsing form: %s", err)
		w.Write([]byte("Bad form data"))
		return
	}

	jsonString := r.URL.Query().Get("json")

	if flagVerbose {
		log.Printf("[Proxy] Got JSON: %s", jsonString)
	}

	var patch StatePatch
	err = json.Unmarshal([]byte(jsonString), &patch)
	if err != nil {
		w.Write([]byte("Bad JSON"))
		return
	}

	select {
	case ph.patchChannel <- patch:
	case <-time.After(time.Second):
		log.Println("[Proxy] Patch channel not receiving patch; discarding")
		w.Write([]byte("Internal error (patch not accepted)"))
	}

	w.Write([]byte("OK"))
}

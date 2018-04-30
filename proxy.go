package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
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
	}
	jsonString := r.PostForm.Get("json")

	if flagVerbose {
		log.Printf("[Proxy] Got JSON: %s", jsonString)
	}

	var patch StatePatch
	json.Unmarshal([]byte(jsonString), &patch)

	ph.patchChannel <- patch

	w.Write([]byte("OK"))
}

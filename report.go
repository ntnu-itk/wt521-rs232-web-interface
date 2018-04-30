package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
)

var flagReportTo string

func init() {
	flag.StringVar(&flagReportTo, "report-to", "", "URL to report new states to (e.g. http://localhost:8080). Empty value => no reporting")
}

func ReportTo(stateChannel <-chan State) {
	for {
		state := <-stateChannel

		data := url.Values{}
		jsonBytes, err := json.Marshal(state)
		data.Set("json", string(jsonBytes))

		log.Printf("[ReportTo] Reporting state %v (json: %s)", state, string(jsonBytes))
		reponse, err := http.PostForm(flagReportTo+"/report/json", data)
		if err != nil {
			log.Printf("[ReportTo] Error while reporting new state: %s", err)
		}
		log.Printf("[ReportTo] Report's response status: %s", reponse.Status)
	}
}

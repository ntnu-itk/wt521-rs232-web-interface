package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"time"
)

var flagReportTo string

func init() {
	flag.StringVar(&flagReportTo, "report-to", "", "URL to report new states to (e.g. http://localhost:8080). Empty value => no reporting")
}

func ReportTo(stateChannel <-chan State) {
	client := &http.Client{Timeout: 10 * time.Second}

	u, err := url.Parse(flagReportTo + "/report/json")
	if err != nil {
		log.Fatalf("[ReportTo] Could not parse URL. Error was: %s", err)
	}

	for {
		state := <-stateChannel

		data := url.Values{}
		jsonBytes, err := json.Marshal(state)
		data.Set("json", string(jsonBytes))
		query := u.Query()
		query.Set("json", string(jsonBytes))
		u.RawQuery = query.Encode()

		if flagVerbose {
			log.Printf("[ReportTo] Reporting state %v (json: %s)", state, string(jsonBytes))
			log.Printf("[ReportTo] URL is %s", u.String())
			log.Printf("[ReportTo] Query: %s", u.RawQuery)
		}

		response, err := client.Get(u.String())
		if err != nil {
			log.Printf("[ReportTo] Error while reporting new state: %s", err)
			continue
		}
		log.Printf("[ReportTo] Report's response status: %s", response.Status)
	}
}

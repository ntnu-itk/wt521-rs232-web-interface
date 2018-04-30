package main

import (
	"flag"
	"log"
)

var flagEnableSerial bool
var flagEnableReporting bool

func main() {
	flag.Parse()
	flagEnableSerial = flagDevice != ""
	flagEnableReporting = flagReportTo != ""

	patchChannel := make(chan StatePatch, 0)
	stateRequestChannel := make(chan *StateRequest, 0)
	currentStateChannel := make(chan State, 0)

	go StateKeeper(patchChannel, stateRequestChannel, currentStateChannel)

	stateHistory := NewStateHistory()
	go stateHistory.Maintain(currentStateChannel)

	if flagEnableSerial {
		if flagVerbose {
			log.Println("[main] flagEnableSerial")
		}
		bytesChannel := make(chan byte, 0)
		go SerialReader(openSerialPort(), bytesChannel)

		mwvMessageChannel := make(chan MWVMessage, 0)
		go MWVMessageConinuousScan(bytesChannel, mwvMessageChannel)

		go MessageToPatchConverter(mwvMessageChannel, patchChannel)
	}

	if flagEnableProxy {
		if flagVerbose {
			log.Println("[main] flagEnableProxy")
		}
		ConfigureProxy(patchChannel)
	}

	if flagReportTo != "" {
		if flagVerbose {
			log.Println("[main] flagReportTo")
		}
		go ReportTo(currentStateChannel)
	} else {
		log.Println(flagReportTo)
	}

	err := HttpServer(stateRequestChannel, currentStateChannel, stateHistory)

	reportError(err)

	log.Fatalf("Error from [HttpServer]: %s", err)
}

// Implement this if you want to report the error, e.g. send an email
func reportError(err error) {}

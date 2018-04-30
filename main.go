package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()

	patchChannel := make(chan StatePatch, 0)
	stateRequestChannel := make(chan *StateRequest, 0)
	currentStateChannel := make(chan State, 0)

	go StateKeeper(patchChannel, stateRequestChannel, currentStateChannel)

	stateHistory := NewStateHistory()
	go stateHistory.Maintain(currentStateChannel)

	if flagEnableSerial {
		bytesChannel := make(chan byte, 0)
		go SerialReader(openSerialPort(), bytesChannel)

		mwvMessageChannel := make(chan MWVMessage, 0)
		go MWVMessageConinuousScan(bytesChannel, mwvMessageChannel)

		go MessageToPatchConverter(mwvMessageChannel, patchChannel)
	}

	if flagEnableProxy {
		log.Println("flagEnableProxy")
		ConfigureProxy(patchChannel)
	}

	if flagReportTo != "" {
		log.Println("flagEnableReporting")
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

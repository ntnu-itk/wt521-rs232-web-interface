package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()

	bytesChannel := make(chan byte, 0)
	go SerialReader(openSerialPort(), bytesChannel)

	mwvMessageChannel := make(chan MWVMessage, 0)
	go MWVMessageConinuousScan(bytesChannel, mwvMessageChannel)

	patchChannel := make(chan StatePatch, 0)
	go MessageToPatchConverter(mwvMessageChannel, patchChannel)

	newStateSubscriptionChannel := make(chan *StateSubscription, 0)
	stateSubscriptionLedger := NewStateSubscriptionLedger()
	go stateSubscriptionLedger.Manage(newStateSubscriptionChannel)

	stateRequestChannel := make(chan StateRequest, 0)
	currentStateChannel := make(chan State, 0)
	go StateKeeper(patchChannel, stateRequestChannel, currentStateChannel, stateSubscriptionLedger)

	log.Fatal(HttpServer(stateRequestChannel, currentStateChannel, newStateSubscriptionChannel))
}

package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()

	patchChannel := make(chan StatePatch, 0)
	stateChannel := make(chan State, 0)
	requestChannel := make(RequestChannel, 0)

	newStateBroadcastChannel := make(chan State, 0)
	newStateReceivedByAllChannel := make(ReceivedByAllChannel, 0)

	stateHistory := NewStateHistory()

	serialPort := openSerialPort()

	mwvByteChannel := make(chan byte, 0)
	mwvMessageChannel := make(chan MWVMessage, 0)

	byteSubscribers := []chan byte{mwvByteChannel}

	stateToBeBroadcastChannel := make(chan State, 0)
	stateSubscriptionChannel := make(chan *StateSubscription, 0)

	go stateHistory.Maintain(NewStateSubscription(Permanent))
	go StateKeeper(patchChannel, requestChannel, stateChannel, stateToBeBroadcastChannel)
	go MessageToPatchConverter(mwvMessageChannel, patchChannel)
	go MWVMessageConinuousScan(mwvByteChannel, mwvMessageChannel)
	go SerialMonitor(serialPort, byteSubscribers)

	log.Fatal(HttpServer(requestChannel, stateChannel, stateSubscriptionChannel, stateHistory))
}

package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()

	patchChannel := make(PatchChannel, 1)
	stateChannel := make(StateChannel, 0)
	requestChannel := make(RequestChannel, 0)

	newStateBroadcastChannel := make(StateChannel, 0)
	newStateReceivedByAllChannel := make(ReceivedByAllChannel, 0)

	stateHistory := NewStateHistory()

	serialPort := openSerialPort()

	go stateHistory.Maintain(newStateBroadcastChannel, newStateReceivedByAllChannel)
	go StateKeeper(patchChannel, requestChannel, stateChannel, newStateBroadcastChannel, newStateReceivedByAllChannel)
	go SerialMonitor(serialPort, patchChannel)

	log.Fatal(HttpServer(requestChannel, stateChannel, newStateBroadcastChannel, newStateReceivedByAllChannel, stateHistory))
}

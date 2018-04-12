package main

import (
	"flag"
	"log"
)

func main() {
	patchChannel := make(PatchChannel, 1)
	stateChannel := make(StateChannel, 0)
	requestChannel := make(RequestChannel, 0)

	flag.Parse()

	serialPort := openSerialPort()

	go stateKeeper(patchChannel, requestChannel, stateChannel)
	go serialMonitor(serialPort, patchChannel)

	log.Fatal(httpServer(requestChannel, stateChannel))
}

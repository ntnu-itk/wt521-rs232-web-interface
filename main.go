package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/tarm/serial"
)

type State struct {
	data string
}

type StatePatch struct {
	newState State
}

type PatchChannel chan StatePatch
type StateChannel chan State
type RequestChannel chan bool

var patchChannel PatchChannel
var stateChannel StateChannel
var requestChannel RequestChannel

func main() {
	patchChannel = make(PatchChannel, 1)
	stateChannel = make(StateChannel, 0)
	requestChannel = make(RequestChannel, 0)

	go stateKeeper(patchChannel, requestChannel, stateChannel)
	go serialMonitor(patchChannel)

	log.Fatal(httpServer(requestChannel, stateChannel))
}

func serialMonitor(patchChannel PatchChannel) {
	var c *serial.Config
	var s *serial.Port
	var err error
	buf := make([]byte, 50)

	c = &serial.Config{
		Name:   "/dev/ttyS0",
		Baud:   300,
		Parity: serial.ParityEven,
		//StopBits: serial.Stop1,
		Size: 7}
	s, err = serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	for true {
		charBuf := make([]byte, 1)
		for true {
			//n, _ :=
			s.Read(charBuf)
			if string(charBuf) == "$" {
				break
			}
		}

		bytesRead := 0

		for bytesRead < 31 {
			n, err := s.Read(buf[bytesRead:])
			bytesRead += n
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Printf("%d*c: %s\n", bytesRead, buf)
		select {
		case <-patchChannel:
			log.Println("Discarded a patch")
		default:
			log.Println("No patch in queue")
		}
		log.Println("Sending patch")
		patchChannel <- StatePatch{
			newState: State{
				data: string(buf)}}
		log.Println("Sent patch")
	}
}

func stateKeeper(patchChannel PatchChannel, requestChannel RequestChannel, stateChannel StateChannel) {
	var state State
	var patch StatePatch

	for {
		select {
		case patch = <-patchChannel:
			log.Println("Got patch")
			state = patch.newState
			log.Printf("New state: %s", state.data)
		case <-requestChannel:
			log.Println("Got request")
			stateChannel <- state
			log.Println("Sent state")
		}
	}
}

func httpServer(requestChannel RequestChannel, stateChannel StateChannel) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request for %s\n", r.URL.Path)
		requestChannel <- true
		log.Println("Sent request")
		state := <-stateChannel
		log.Println("Got state")
		fmt.Fprintf(w, "<!DOCTYPE html><html><body>Hello, %q\n<br>\nState: %s", html.EscapeString(r.URL.Path), state.data)
	})
	return http.ListenAndServe(":8080", nil)
}

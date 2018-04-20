package main

import (
	"container/list"
	"flag"
	"log"
	"time"
)

var flagHistoryLimit int
var flagMaxSubscribers int

func init() {
	flag.IntVar(&flagHistoryLimit, "history", 1000, "max number of readings to keep in memory")
	flag.IntVar(&flagMaxSubscribers, "max-subscribers", 1000, "max number of requests that can wait for new data; exceeding it will discard oldest request")
}

func StateKeeper(
	patchChannel chan StatePatch,
	requestChannel RequestChannel,
	stateChannel chan State,
	stateToBeBroadcastChannel chan<- State) {
	state := State{}
	patch := StatePatch{}

	for {
		select {
		case patch = <-patchChannel:
			state.Patch(patch)
			stateToBeBroadcastChannel <- state
		case <-requestChannel:
			stateChannel <- state
		}
	}
}

func (state *State) SendToAll(stateChannel chan State, receivedByAllChannel ReceivedByAllChannel) (nSends int) {
	if flagVerbose {
		log.Println("Sending new State to all listeners…")
	}

SendToAllLoop:
	for nSends = 0; ; nSends++ {
		select {
		case stateChannel <- *state:
			// OK
		default:
			break SendToAllLoop
		}
	}

AllReceivedLoop:
	for {
		select {
		case receivedByAllChannel <- true:
			// OK
		default:
			break AllReceivedLoop
		}
	}

	return
}

// Updates the state history when new states occur
func (stateHistory StateHistory) Maintain(newStateChannel chan State, receivedByAllChannel ReceivedByAllChannel) {
	for {
		state := <-newStateChannel
		<-receivedByAllChannel

		stateHistory.list.PushBack(state)

		for stateHistory.list.Len() > stateHistory.maxLength {
			stateHistory.list.Remove(stateHistory.list.Front())
		}
	}
}

func MessageToPatchConverter(messageChannel <-chan MWVMessage, patchChannel chan<- StatePatch) {
	var patch StatePatch
	var message MWVMessage
	for {
		message = <-messageChannel
		if flagVerbose {
			log.Printf("Converting %v to patch…", message)
		}
		patch.windAngle = message.dir
		patch.windSpeed = message.spd
		patch.lastUpdated = time.Now()
		if flagVerbose {
			log.Printf("Sending patch %v on patch channel…", patch)
		}
		patchChannel <- patch
		if flagVerbose {
			log.Println("Patch sent")
		}
	}
}

func StateBroadcaster(newSubcriptionsChannel <-chan *StateSubscription, stateToBeBroadcastChannel <-chan State) {
	subscriberList := list.New()

	for {
		select {
		case state := <-stateToBeBroadcastChannel:
			for e := subscriberList.Front(); e != nil; e = e.Next() {
				e.Value.(chan State) <- state
			}
		case subscriber := <-newSubcriptionsChannel:
			subscriberList.PushBack(subscriber)
			if subscriberList.Len() > flagMaxSubscribers {
				subscriberList.Remove(subscriberList.Front())
			}
		}
	}
}

package main

import "log"

func StateKeeper(
	patchChannel PatchChannel,
	requestChannel RequestChannel,
	stateChannel StateChannel,
	newStateBroadcastChannel StateChannel,
	receivedByAllChannel ReceivedByAllChannel) {
	state := State{}
	patch := StatePatch{}

	for {
		select {
		case patch = <-patchChannel:
			state.Patch(patch)
			nSends := state.SendToAll(newStateBroadcastChannel, receivedByAllChannel)
			if flagVerbose {
				log.Printf("Sent new State to %d listeners.", nSends)
			}
		case <-requestChannel:
			stateChannel <- state
		}
	}
}

func (state *State) SendToAll(stateChannel StateChannel, receivedByAllChannel ReceivedByAllChannel) (nSends int) {
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

package main

func stateKeeper(patchChannel PatchChannel, requestChannel RequestChannel, stateChannel StateChannel) {
	state := State{}
	patch := StatePatch{}

	for {
		select {
		case patch = <-patchChannel:
			state.patch(patch)
		case <-requestChannel:
			stateChannel <- state
		}
	}
}

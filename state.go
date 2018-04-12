package main

import (
	"fmt"
	"log"
	"time"
)

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
		log.Println(&state)
	}
}

func (state *State) patch(patch StatePatch) {
	state.windSpeed = patch.windSpeed
	state.windAngle = patch.windAngle
	state.lastUpdated = time.Now()
}

func (state *State) String() string {
	return fmt.Sprintf("State{ Speed:%f  Angle:%d  Updated:%s }",
		state.windSpeed,
		state.windAngle,
		state.lastUpdated)
}

func (patch *StatePatch) String() string {
	return fmt.Sprintf("StatePatch{ Speed:%f  Angle:%d }",
		patch.windSpeed,
		patch.windAngle)
}

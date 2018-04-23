package main

import (
	"fmt"
	"time"
)

type State struct {
	windAngle   WindAngle
	windSpeed   WindSpeed
	lastUpdated time.Time
}

type StateRequest struct {
	reply chan State
}

func NewStateRequest() *StateRequest {
	return &StateRequest{
		reply: make(chan State, 0)}
}

func StateKeeper(
	patchChannel <-chan StatePatch,
	requestChannel <-chan *StateRequest,
	stateChannel chan<- State) {
	state := State{}

	for {
		select {
		case patch := <-patchChannel:
			state.Apply(patch)
		BroadcastLoop:
			for {
				select {
				case stateChannel <- state:
				default:
					break BroadcastLoop
				}
			}
			break
		case request := <-requestChannel:
			request.reply <- state
			break
		}
	}
}

func (state *State) Apply(patch StatePatch) {
	state.windSpeed = patch.windSpeed
	state.windAngle = patch.windAngle
	state.lastUpdated = time.Now()
}

func (state *State) String() string {
	return fmt.Sprintf("State{Speed:%.1f, Angle:%d, Updated:%s}",
		state.windSpeed,
		state.windAngle,
		SimpleTimeString(state.lastUpdated))
}

func (state *State) ToJSON() string {
	return fmt.Sprintf(`{"speed":%f,"angle":%d,"time":"%s"}`, state.windSpeed, state.windAngle, SimpleTimeString(state.lastUpdated))
}

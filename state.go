package main

import (
	"fmt"
	"time"
)

type State struct {
	WindAngle   WindAngle `json:"angle"`
	WindSpeed   WindSpeed `json:"speed"`
	LastUpdated time.Time `json:"time"`
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
	state.WindSpeed = patch.WindSpeed
	state.WindAngle = patch.WindAngle
	state.LastUpdated = time.Now()
}

func (state *State) ToJSON() string {
	return fmt.Sprintf(`{
    "speed": %.1f,
    "angle": %d,
    "time": "%s"
}`,
		state.WindSpeed,
		state.WindAngle,
		state.LastUpdated.String())
}

func (state *State) String() string {
	return fmt.Sprintf("State{Speed:%.1f, Angle:%d, Updated:%s}",
		state.WindSpeed,
		state.WindAngle,
		state.LastUpdated.String())
}

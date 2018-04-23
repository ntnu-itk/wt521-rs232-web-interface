package main

import (
	"flag"
	"fmt"
	"time"
)

type State struct {
	windAngle   WindAngle
	windSpeed   WindSpeed
	lastUpdated time.Time
}

type StateRequest bool

const GimmePls StateRequest = true

var flagHistoryLimit int
var flagMaxSubscribers int

func init() {
	flag.IntVar(&flagHistoryLimit, "history", 1000, "max number of readings to keep in memory")
	flag.IntVar(&flagMaxSubscribers, "max-subscribers", 1000, "max number of requests that can wait for new data; exceeding it will discard oldest request")
}

func StateKeeper(
	patchChannel <-chan StatePatch,
	requestChannel <-chan StateRequest,
	stateChannel chan<- State,
	stateSubscriptionLedger *StateSubscriptionLedger) {
	state := State{}

	for {
		select {
		case patch := <-patchChannel:
			state.Apply(patch)
			stateSubscriptionLedger.Broadcast(state)
		case <-requestChannel:
			stateChannel <- state
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

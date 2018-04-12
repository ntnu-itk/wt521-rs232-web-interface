package main

import (
	"fmt"
	"time"
)

type State struct {
	windAngle   int16
	windSpeed   float64
	lastUpdated time.Time
}

type StatePatch State

type PatchChannel chan StatePatch
type StateChannel chan State
type RequestChannel chan bool

func (state *State) patch(patch StatePatch) {
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

func (patch *StatePatch) String() string {
	return fmt.Sprintf("StatePatch{Speed:%.1f, Angle:%d }",
		patch.windSpeed,
		patch.windAngle)
}

func SimpleTimeString(t time.Time) string {
	str := t.String()
	if len(str) < 40 {
		return str
	}
	return str[:19] + str[29:40]
}

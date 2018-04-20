package main

import (
	"container/list"
	"fmt"
	"time"
)

type State struct {
	windAngle   WindAngle
	windSpeed   WindSpeed
	lastUpdated time.Time
}

type StatePatch State

type StateHistory struct {
	list      *list.List
	maxLength int
}

type StateSubscription struct {
	ch        chan State
	permanent SubscriptionPersistence
}

type SubscriptionPersistence int

const (
	Permanent SubscriptionPersistence = -1
	OneShot   SubscriptionPersistence = 1
)

func NewStateSubscription(permanent SubscriptionPersistence) *StateSubscription {
	return &StateSubscription{ch: make(chan State, 0), permanent: permanent}
}

type RequestChannel chan bool

type ReceivedByAllChannel chan bool

func (state *State) Patch(patch StatePatch) {
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

func (state *State) ToJSON() string {
	return fmt.Sprintf(`{"speed":%f,"angle":%d,"time":"%s"}`, state.windSpeed, state.windAngle, SimpleTimeString(state.lastUpdated))
}

func (stateHistory *StateHistory) ToJSON() (str string) {
	var state State
	separator := ""
	for e := stateHistory.list.Front(); e != nil; e = e.Next() {
		state = e.Value.(State)
		str += separator + state.ToJSON()
		separator = ","
	}
	return
}

func SimpleTimeString(t time.Time) string {
	var (
		date           string
		time           string
		unusedInt      int
		timeZoneOffset string
		timeZoneName   string
	)

	fmt.Sscanf(t.String(),
		"%10s %8s.%d %s %s",
		&date,
		&time,
		&unusedInt,
		&timeZoneOffset,
		&timeZoneName)

	return fmt.Sprintf("%s %s",
		date, time)
}

func NewStateHistory() *StateHistory {
	return &StateHistory{
		list:      list.New(),
		maxLength: flagHistoryLimit}
}

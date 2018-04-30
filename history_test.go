package main

import (
	"testing"
	"time"
)

/*
type StateHistory struct {
	list      *list.List
	maxLength int
}

func NewStateHistory() *StateHistory {
	return &StateHistory{
		list:      list.New(),
		maxLength: flagHistorySamples}
}

// Updates the state history when new states occur
func (stateHistory *StateHistory) Maintain(stateChannel chan State) {
	for {
		state := <-stateChannel

		stateHistory.list.PushBack(state)

		for stateHistory.list.Len() > stateHistory.maxLength {
			stateHistory.list.Remove(stateHistory.list.Front())
		}
	}
}

func (stateHistory *StateHistory) AsSlice() []State {
	slice := make([]State, stateHistory.list.Len())

	i := 0
	for e := stateHistory.list.Front(); e != nil; e = e.Next() {
		slice[i] = e.Value.(State)
		i++
	}

	return slice
}
*/

func TestAcceptsStates(t *testing.T) {
	history := NewStateHistory()
	stateChannel := make(chan State)

	go history.Maintain(stateChannel)

	for i := 0; i < 4; i++ {
		state := State{
			WindAngle:   WindAngle(2 * i),
			WindSpeed:   WindSpeed(1.2 * float64(i)),
			LastUpdated: time.Now()}

		select {
		case stateChannel <- state:
		case <-time.After(time.Millisecond):
			t.Errorf("does not receive state #%d", i)
		}
	}
}

func TestAsSlice(t *testing.T) {
	history := NewStateHistory()
	stateChannel := make(chan State)

	go history.Maintain(stateChannel)

	stateChannel <- State{
		WindAngle:   2,
		WindSpeed:   3.4,
		LastUpdated: time.Now()}
	stateChannel <- State{
		WindAngle:   4,
		WindSpeed:   5.6,
		LastUpdated: time.Now()}

	<-time.After(time.Millisecond)
	slice := history.AsSlice()

	if len(slice) != 2 {
		t.Errorf("len should be 2, was %d:\n\t%v", len(slice), slice)
	}

	stateChannel <- State{
		WindAngle:   6,
		WindSpeed:   7.8,
		LastUpdated: time.Now()}

	<-time.After(time.Millisecond)
	slice = history.AsSlice()

	if len(slice) != 3 {
		t.Errorf("len should be 3, was %d:\n\t%v", len(slice), slice)
	}
}

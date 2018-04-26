package main

import "container/list"

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

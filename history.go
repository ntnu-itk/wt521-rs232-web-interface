package main

import (
	"container/list"
	"flag"
)

var flagHistoryLimit int

func init() {
	flag.IntVar(&flagHistoryLimit, "history", 1000, "max number of readings to keep in memory")
}

type StateHistory struct {
	list      *list.List
	maxLength int
}

func NewStateHistory() *StateHistory {
	return &StateHistory{
		list:      list.New(),
		maxLength: flagHistoryLimit}
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

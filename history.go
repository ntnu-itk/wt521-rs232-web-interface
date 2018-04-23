package main

import "container/list"

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
func (stateHistory *StateHistory) Maintain(newSubscriptionChannel chan<- *StateSubscription) {
	ss := NewStateSubscription(Permanent)
	for {
		state := <-ss.ch

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

package main

import (
	"testing"
	"time"
)

func TestStateKeeperKeepsAndReplies(t *testing.T) {
	patchChannel := make(chan StatePatch)
	requestChannel := make(chan *StateRequest)
	stateChannel := make(chan State)

	go StateKeeper(patchChannel, requestChannel, stateChannel)

	patch := StatePatch{
		WindAngle: 1,
		WindSpeed: 2.3}
	select {
	case patchChannel <- patch:
	case <-time.After(time.Millisecond):
		t.Error("does not receive patch")
	}

	request := NewStateRequest()
	select {
	case requestChannel <- request:
	case <-time.After(time.Millisecond):
		t.Error("does not receive request")
	}

	var state State
	select {
	case state = <-request.reply:
	case <-time.After(time.Millisecond):
		t.Error("does not reply to request")
	}
	if state.WindAngle != patch.WindAngle || state.WindSpeed != patch.WindSpeed {
		t.Errorf("state was not patched properly. State: %v, Patch: %v", state, patch)
	}
}

func TestStateKeeperKeepsAndBroadcasts(t *testing.T) {
	patchChannel := make(chan StatePatch)
	requestChannel := make(chan *StateRequest)
	stateChannel := make(chan State)

	go StateKeeper(patchChannel, requestChannel, stateChannel)

	patch := StatePatch{
		WindAngle: 1,
		WindSpeed: 2.3}

	for i := 0; i < 5; i++ {
		go func(stateChannel chan State, id int, t *testing.T, patch StatePatch) {
			select {
			case state := <-stateChannel:
				if state.WindAngle != patch.WindAngle || state.WindSpeed != patch.WindSpeed {
					t.Errorf("state was not patched properly. State: %v, Patch: %v", state, patch)
				}
				break
			case <-time.After(time.Millisecond):
				t.Errorf("broadcast was not received by listener #%d", id)
			}
		}(stateChannel, i, t, patch)
	}

	select {
	case patchChannel <- patch:
	case <-time.After(time.Millisecond):
		t.Error("does not receive patch")
	}
}

func TestStateApplyPatch(t *testing.T) {
	originalState := State{
		WindAngle:   42,
		WindSpeed:   13.37,
		LastUpdated: time.Now()}

	_bogusTime, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2006")
	bogusTime := _bogusTime
	patch := StatePatch{
		WindAngle:   64,
		WindSpeed:   7.2,
		LastUpdated: bogusTime}

	state := originalState

	<-time.After(1 * time.Microsecond)

	state.Apply(patch)
	if !state.LastUpdated.After(originalState.LastUpdated) {
		t.Errorf("time was not updated")
	}
	if state.WindAngle != patch.WindAngle {
		t.Error("wind angle not updated")
	}
	if state.WindSpeed != patch.WindSpeed {
		t.Error("wind angle not updated")
	}
	if state.LastUpdated.Unix() == bogusTime.Unix() {
		t.Error("patching should not use the patch's time")
	}
}

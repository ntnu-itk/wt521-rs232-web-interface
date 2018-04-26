package main

import (
	"testing"
	"time"
)

func TestMessageToPatchConverter(t *testing.T) {
	messageChannel := make(chan MWVMessage)
	patchChannel := make(chan StatePatch)

	go MessageToPatchConverter(messageChannel, patchChannel)

	var message MWVMessage

	expectedAngle := WindAngle(170)
	expectedSpeed := WindSpeed(4.8)
	err := message.parse("$WIMWV,170,R,4.8,M,A*34")
	if err != nil {
		t.Errorf("test code may be bad; failed to parse litteral message string with error: %s", err)
		return
	}

	select {
	case messageChannel <- message:
	case <-time.After(time.Millisecond):
		t.Error("does not receive message")
	}

	select {
	case patch := <-patchChannel:
		if patch.WindAngle != expectedAngle || patch.WindSpeed != expectedSpeed {
			t.Errorf("bad patch; expected angle %d and speed %f but got:\n\t%v", expectedAngle, expectedSpeed, patch)
		}
	}
}

func TestStatePatchString(t *testing.T) {
	patch := StatePatch{
		WindAngle: 123,
		WindSpeed: 4.56}
	expected := "StatePatch{Speed:4.6, Angle:123}"
	if patch.String() != expected {
		t.Errorf("expected\n\t%s\ngot\n\t%s", expected, patch.String())
	}
}

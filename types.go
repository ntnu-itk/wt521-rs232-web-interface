package main

import "time"

type State struct {
	windAngle   int16
	windSpeed   float64
	lastUpdated time.Time
}

type StatePatch State

type PatchChannel chan StatePatch
type StateChannel chan State
type RequestChannel chan bool

package main

import (
	"fmt"
	"time"
)

type MyTime time.Time

func (mt *MyTime) Time() time.Time {
	return time.Time(*mt)
}

func (t *MyTime) String() string {
	var (
		date           string
		clockTime      string
		unusedInt      int
		timeZoneOffset string
		timeZoneName   string
	)

	fmt.Sscanf(t.Time().String(),
		"%10s %8s.%d %s %s",
		&date,
		&clockTime,
		&unusedInt,
		&timeZoneOffset,
		&timeZoneName)

	return fmt.Sprintf("%s %s",
		date, clockTime)
}

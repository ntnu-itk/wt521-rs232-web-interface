package main

import (
	"fmt"
	"time"
)

type MyTime time.Time

func (t *MyTime) String() string {
	var (
		date           string
		clockTime      string
		unusedInt      int
		timeZoneOffset string
		timeZoneName   string
	)

	fmt.Sscanf(time.Time(*t).String(),
		"%10s %8s.%d %s %s",
		&date,
		&clockTime,
		&unusedInt,
		&timeZoneOffset,
		&timeZoneName)

	return fmt.Sprintf("%s %s",
		date, clockTime)
}

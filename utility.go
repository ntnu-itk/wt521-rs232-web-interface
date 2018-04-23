package main

import (
	"fmt"
	"time"
)

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

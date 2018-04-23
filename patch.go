package main

import (
	"fmt"
	"log"
	"time"
)

type StatePatch State

func (patch *StatePatch) String() string {
	return fmt.Sprintf("StatePatch{Speed:%.1f, Angle:%d }",
		patch.windSpeed,
		patch.windAngle)
}

func MessageToPatchConverter(messageChannel <-chan MWVMessage, patchChannel chan<- StatePatch) {
	var patch StatePatch
	for {
		message := <-messageChannel
		if flagVerbose {
			log.Printf("Converting %v to patch…", message)
		}
		patch.windAngle = message.dir
		patch.windSpeed = message.spd
		patch.lastUpdated = time.Now()
		if flagVerbose {
			log.Printf("Sending patch %v on patch channel…", patch)
		}
		patchChannel <- patch
		if flagVerbose {
			log.Println("Patch sent")
		}
	}
}

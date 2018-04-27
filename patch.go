package main

import (
	"fmt"
	"log"
	"time"
)

type StatePatch State

func (patch *StatePatch) String() string {
	return fmt.Sprintf("StatePatch{Speed:%.1f, Angle:%d}",
		patch.WindSpeed,
		patch.WindAngle)
}

func MessageToPatchConverter(messageChannel <-chan MWVMessage, patchChannel chan<- StatePatch) {
	var patch StatePatch
	for {
		message := <-messageChannel

		if flagVerbose {
			log.Printf("[MessageToPatchConverter] Converting %v to patch…", message)
		}

		patch.WindAngle = message.dir
		patch.WindSpeed = message.spd
		patch.LastUpdated = MyTime(time.Now())

		if flagVerbose {
			log.Printf("[MessageToPatchConverter] Sending patch %v on patch channel…", patch)
		}

		patchChannel <- patch

		if flagVerbose {
			log.Println("[MessageToPatchConverter] Patch sent")
		}
	}
}

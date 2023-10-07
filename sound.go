package main

import (
	"github.com/gen2brain/beeep"
)

func playBeep(beeps int) {
	for i := 0; i < beeps; i++ {
		beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
	}
}

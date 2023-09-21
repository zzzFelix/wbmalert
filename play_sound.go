package main

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var soundInitialized = false

func playSound() {
	f, err := os.Open("success.mp3")
	if err != nil {
		fmt.Println(err)
	}

	streamer, format, err := mp3.Decode(f)

	defer streamer.Close()
	if !soundInitialized {
		err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			fmt.Println(err)
		}
		soundInitialized = true
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}

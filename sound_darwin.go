//go:build darwin

package main

import (
	"embed"
	"log"
	"os/exec"
)

//go:embed success.mp3
var staticAssets embed.FS

func playSound() {
	const soundFileName = "success.mp3"

	f, err := staticAssets.Open(soundFileName)
	if err != nil {
		log.Println(err)
		return
	}

	cmd := exec.Command("afplay", soundFileName)

	err = cmd.Run()
	if err != nil {
		log.Println("Error playing sound:", err)
	}
	f.Close()
}

//go:build !darwin

package main

import (
	"log"
	"runtime"
)

func playSound() {
	log.Println("Sound not supported for OS %s", runtime.GOOS)
}

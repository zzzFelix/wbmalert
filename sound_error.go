//go:build !darwin

package main

import (
	"log"
	"runtime"
)

func playSound() {
	log.Printf("Sound not supported for OS %s", runtime.GOOS)
}

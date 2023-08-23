package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func main() {
	for {
		checkWBM()
		time.Sleep(240 * time.Second)
	}
}

func checkWBM() {
	resp, err := http.Get("https://www.wbm.de/wohnungen-berlin/angebote/")
	if err != nil {
		fmt.Println("Es ist ein Fehler aufgetreten")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	content := string(body[:])
	if reference := "LEIDER HABEN WIR DERZEIT KEINE VERFÜGBAREN WOHNUNGSANGEBOTE"; strings.Contains(content, reference) {
		fmt.Println("Keine Ergebnisse")
	} else if strings.Contains(content, "403 Zugriff verweigert") {
		fmt.Println("403 Zugriff verweigert")
	} else {
		fmt.Println("Neue Ergebnisse!!!!")
		playSound()
	}
}

func playSound() {
	f, err := os.Open("success.mp3")
	if err != nil {
		fmt.Println(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Println(err)
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}

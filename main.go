package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/microcosm-cc/bluemonday"
)

var interval = 90 // in seconds
var websites = []Website{
	{"WBM", "https://www.wbm.de/wohnungen-berlin/angebote/", ""},
	{"Gewobag", "https://www.gewobag.de/fuer-mieter-und-mietinteressenten/mietangebote/?bezirke%5B%5D=charlottenburg-wilmersdorf-charlottenburg&bezirke%5B%5D=friedrichshain-kreuzberg&bezirke%5B%5D=friedrichshain-kreuzberg-friedrichshain&bezirke%5B%5D=friedrichshain-kreuzberg-kreuzberg&bezirke%5B%5D=pankow&bezirke%5B%5D=pankow-prenzlauer-berg&bezirke%5B%5D=tempelhof-schoeneberg-schoeneberg&objekttyp%5B%5D=wohnung&gesamtmiete_von=&gesamtmiete_bis=&gesamtflaeche_von=&gesamtflaeche_bis=&zimmer_von=&zimmer_bis=&keinwbs=1&sort-by=recent/", ""},
}

type Website struct {
	name     string
	url      string
	snapshot string
}

func main() {
	initializeWebsites()

	for {
		for i := 0; i < len(websites); i++ {
			websites[i] = checkWebsite(websites[i])
		}
		goToSleep()
	}
}

func initializeWebsites() {
	for i := 0; i < len(websites); i++ {
		websites[i] = createInitialSnapshot(websites[i])
		fmt.Println("Created initial snapshot for " + websites[i].name)
	}
	goToSleep()
}

func createInitialSnapshot(website Website) Website {
	content := getWebsiteAsString(website)
	website.snapshot = content
	return website
}

func getWebsiteAsString(website Website) string {
	resp, err := http.Get(website.url)
	if err != nil {
		fmt.Println(err)
		fmt.Println("An error occurred! The website could not be reached!")
		return "Error"
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	content := string(body[:])
	content = sanitizeHtml(content)
	content = removeAllWhitespace(content)
	return content
}

func removeAllWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func checkWebsite(website Website) Website {
	content := getWebsiteAsString(website)

	if website.snapshot != content {
		fmt.Println("========= " + website.name + " =========")
		fmt.Println("Content changed: " + website.url)
		fmt.Println("====================" + strings.Repeat("=", len(website.name)))
		playSound()
		website.snapshot = content
	} else {
		fmt.Println("No changes for " + website.name)
	}

	return website
}

func goToSleep() {
	fmt.Printf("Going to sleep for %s seconds", strconv.FormatInt(int64(interval), 10))
	fmt.Println()
	time.Sleep(time.Duration(interval) * time.Second)
}

func sanitizeHtml(s string) string {
	p := bluemonday.StripTagsPolicy()
	html := p.Sanitize(s)
	return html
}

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

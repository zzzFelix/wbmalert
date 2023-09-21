package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type website struct {
	Name     string
	Url      string
	Snapshot string
}

var interval = 0 // in seconds
var websites = []website{}

func main() {
	configuration := readConfiguration()
	initializeWebsites(configuration)
	goToSleep()

	for {
		for i := 0; i < len(websites); i++ {
			checkWebsite(websites[i])
		}
		goToSleep()
	}
}

func initializeWebsites(configuration configuration) {
	websites = configuration.Websites
	interval = configuration.Interval

	for i := 0; i < len(websites); i++ {
		createInitialSnapshot(&websites[i])
	}
}

func createInitialSnapshot(website *website) {
	content, error := getWebsiteAsString(website)
	if error == nil {
		website.Snapshot = content
		log.Println("Created initial snapshot for " + website.Name)
	}
}

func getWebsiteAsString(website *website) (string, error) {
	resp, err := http.Get(website.Url)
	if err != nil {
		log.Println(err)
		return "Error", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	content := string(body[:])
	content = sanitizeHtml(content)
	return content, nil
}

func checkWebsite(website website) {
	content, error := getWebsiteAsString(&website)
	if error != nil {
		return
	}
	if website.Snapshot != content {
		website.Snapshot = content
		printContentChangeMsg(website)
		playSound()
	} else {
		log.Println("No changes for " + website.Name)
	}
}

func printContentChangeMsg(website website) {
	log.Println("========= " + website.Name + " =========")
	log.Println("Content changed: " + website.Url)
	log.Println("====================" + strings.Repeat("=", len(website.Name)))
}

func goToSleep() {
	log.Printf("Going to sleep for %s seconds", strconv.FormatInt(int64(interval), 10))
	time.Sleep(time.Duration(interval) * time.Second)
}

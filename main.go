package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Website struct {
	Name     string
	Url      string
	Snapshot string
}

var interval = 0 // in seconds
var websites = []Website{}

func main() {
	configuration := readConfiguration()
	initializeWebsites(configuration)
	goToSleep()

	for {
		for i := 0; i < len(websites); i++ {
			websites[i] = checkWebsite(websites[i])
		}
		goToSleep()
	}
}

func initializeWebsites(configuration Configuration) {
	websites = configuration.Websites
	interval = configuration.Interval

	for i := 0; i < len(websites); i++ {
		createInitialSnapshot(websites[i])
	}
}

func createInitialSnapshot(website Website) {
	content, error := getWebsiteAsString(website)
	if error == nil {
		website.Snapshot = content
	}
	fmt.Println("Created initial snapshot for " + website.Name)
}

func getWebsiteAsString(website Website) (string, error) {
	resp, err := http.Get(website.Url)
	if err != nil {
		fmt.Println(err)
		fmt.Println("An error occurred! The website could not be reached!")
		return "Error", errors.New("The website could not be reached")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	content := string(body[:])
	content = sanitizeHtml(content)
	return content, nil
}

func checkWebsite(website Website) Website {
	content, error := getWebsiteAsString(website)
	if error != nil {
		return website
	}
	if website.Snapshot != content {
		website.Snapshot = content
		printContentChangeMsg(website)
		playSound()
	} else {
		fmt.Println("No changes for " + website.Name)
	}

	return website
}

func printContentChangeMsg(website Website) {
	fmt.Println("========= " + website.Name + " =========")
	fmt.Println("Content changed: " + website.Url)
	fmt.Println("====================" + strings.Repeat("=", len(website.Name)))
}

func goToSleep() {
	fmt.Printf("Going to sleep for %s seconds", strconv.FormatInt(int64(interval), 10))
	fmt.Println()
	time.Sleep(time.Duration(interval) * time.Second)
}

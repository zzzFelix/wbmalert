package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	client   httpClient
	interval int
	websites []website
)

func init() {
	client = &http.Client{}
	interval = 0
	websites = []website{}
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type website struct {
	Name         string
	Url          string
	RegexpRemove string
	Snapshot     string
}

func main() {
	configuration := readConfiguration()
	initializeWebsites(configuration)

	for {
		var wg sync.WaitGroup
		for i := range websites {
			wg.Add(1)
			go checkWebsite(&websites[i], &wg)
		}

		wg.Wait()
		goToSleep()
	}
}

func initializeWebsites(configuration configuration) {
	var wg sync.WaitGroup
	websites = configuration.Websites
	interval = configuration.Interval

	for i := range websites {
		wg.Add(1)
		go createInitialSnapshot(&websites[i], &wg)
	}

	wg.Wait()
	goToSleep()
}

func createInitialSnapshot(website *website, wg *sync.WaitGroup) {
	content, err := getWebsiteAsString(website)
	if err == nil {
		website.Snapshot = content
		log.Println("Created initial snapshot for " + website.Name)
	} else {
		log.Println(err)
	}
	defer wg.Done()
}

func getWebsiteAsString(website *website) (string, error) {
	request, err := http.NewRequest(http.MethodGet, website.Url, nil)
	if err != nil {
		return "Error", err
	}
	resp, err := client.Do(request)
	if err != nil {
		return "Error", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error", err
	}
	content := string(body[:])
	content = removeByRegexp(content, website.RegexpRemove)
	content = sanitizeHtml(content)
	return content, nil
}

func checkWebsite(website *website, wg *sync.WaitGroup) {
	content, err := getWebsiteAsString(website)

	if err == nil {
		if website.Snapshot != content {
			website.Snapshot = content
			printContentChangeMsg(website)
			playSound()
		} else {
			log.Println("No changes for " + website.Name)
		}
	} else {
		log.Println(err)
	}

	defer wg.Done()
}

func printContentChangeMsg(website *website) {
	log.Println("========= " + website.Name + " =========")
	log.Println("Content changed: " + website.Url)
	log.Println("====================" + strings.Repeat("=", len(website.Name)))
}

func goToSleep() {
	log.Printf("Going to sleep for %d seconds", interval)
	time.Sleep(time.Duration(interval) * time.Second)
}

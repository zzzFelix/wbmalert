package main

import (
	"fmt"
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
		log.Println("[INFO]", "Created initial snapshot for "+website.Name)
	} else {
		log.Println("[ERROR]", website.Name, err)
	}
	defer wg.Done()
}

func getWebsiteAsString(website *website) (string, error) {
	request, err := http.NewRequest(http.MethodGet, website.Url, nil)
	errorText := fmt.Sprintf("Error fetching %s", website.Name)
	if err != nil {
		return errorText, err
	}
	resp, err := client.Do(request)
	if err != nil {
		return errorText, err
	}
	if resp.StatusCode >= 400 {
		return errorText, fmt.Errorf("http status code error occured : %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errorText, err
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
			log.Println("[INFO]", "No changes for "+website.Name)
		}
	} else {
		log.Println("[ERROR]", website.Name, err)
	}

	defer wg.Done()
}

func printContentChangeMsg(website *website) {
	log.Println("[INFO]", "========= "+website.Name+" =========")
	log.Println("[INFO]", "Content changed: "+website.Url)
	log.Println("[INFO]", "===================="+strings.Repeat("=", len(website.Name)))
}

func goToSleep() {
	log.Printf("[INFO] Going to sleep for %d seconds", interval)
	time.Sleep(time.Duration(interval) * time.Second)
}

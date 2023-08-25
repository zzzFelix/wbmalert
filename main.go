package main

import (
	"errors"
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
	//{"WBM", "https://www.wbm.de/wohnungen-berlin/angebote/", ""},
	{"Gesobau", "https://www.gesobau.de/mieten/wohnungssuche/", ""},
	{"Stadt und Land", "https://www.stadtundland.de/immobiliensuche.php?form=stadtundland-expose-search-1.form&sp%3Acategories%5B3352%5D%5B%5D=-&sp%3Acategories%5B3352%5D%5B%5D=__last__&sp%3AroomsFrom%5B%5D=&sp%3AroomsTo%5B%5D=&sp%3ArentPriceFrom%5B%5D=&sp%3ArentPriceTo%5B%5D=&sp%3AareaFrom%5B%5D=&sp%3AareaTo%5B%5D=&sp%3Afeature%5B%5D=__last__&action=submit", ""},
	{"Degewo", "https://immosuche.degewo.de/de/search?size=10&page=1&property_type_id=1&categories%5B%5D=1&lat=&lon=&area=&address%5Bstreet%5D=&address%5Bcity%5D=&address%5Bzipcode%5D=&address%5Bdistrict%5D=&address%5Braw%5D=&district=&property_number=&price_switch=true&price_radio=null&price_from=&price_to=&qm_radio=null&qm_from=&qm_to=&rooms_radio=null&rooms_from=&rooms_to=&wbs_required=false&order=rent_total_without_vat_asc", ""},
	{"Gewobag", "https://www.gewobag.de/fuer-mieter-und-mietinteressenten/mietangebote/?bezirke%5B%5D=charlottenburg-wilmersdorf-charlottenburg&bezirke%5B%5D=friedrichshain-kreuzberg&bezirke%5B%5D=friedrichshain-kreuzberg-friedrichshain&bezirke%5B%5D=friedrichshain-kreuzberg-kreuzberg&bezirke%5B%5D=pankow&bezirke%5B%5D=pankow-prenzlauer-berg&bezirke%5B%5D=tempelhof-schoeneberg-schoeneberg&objekttyp%5B%5D=wohnung&gesamtmiete_von=&gesamtmiete_bis=&gesamtflaeche_von=&gesamtflaeche_bis=&zimmer_von=&zimmer_bis=&keinwbs=1&sort-by=recent/", ""},
	{"Core Immobilienmanagement", "https://www.core-berlin.de/de/vermietung", ""},
	{"Heimstaden", "https://portal.immobilienscout24.de/ergebnisliste/92830022", ""},
	{"DPF", "https://www.dpfonline.de/interessenten/immobilien/", ""},
	{"ImmoScout", "https://www.immobilienscout24.de/Suche/de/berlin/berlin/wohnung-mieten?numberofrooms=2.0-&price=-1300.0&exclusioncriteria=swapflat&pricetype=calculatedtotalrent&geocodes=110000000104,110000000801,110000000101,110000000102,110000000106,110000000301,110000000901,110000000201,110000000202,110000001103&enteredFrom=saved_search&utm_medium=email&utm_source=sfmc&utm_campaign=fulfillment_update&utm_content=fulfillment_results", ""},
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
	content, error := getWebsiteAsString(website)
	if error != nil {
		website.snapshot = content
	}
	return website
}

func getWebsiteAsString(website Website) (string, error) {
	resp, err := http.Get(website.url)
	if err != nil {
		fmt.Println(err)
		fmt.Println("An error occurred! The website could not be reached!")
		return "Error", errors.New("The website could not be reached")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	content := string(body[:])
	content = sanitizeHtml(content)
	content = removeAllWhitespace(content)
	return content, nil
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
	content, error := getWebsiteAsString(website)
	if error != nil {
		return website
	}
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

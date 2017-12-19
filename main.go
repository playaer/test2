package main

import (
	"log"
	"net/http"
	"sync"
	"strings"
	"regexp"
	"io/ioutil"
	"encoding/json"
)

type SiteMeta struct {
	Status  int                 `json:"status"`
	Headers []map[string]string `json:"headers"`
}

type SiteElements struct {
	TagName string `json:"tag-name"`
	Count   int    `json:"count"`
}

type SiteData struct {
	Url      string          `json:"url"`
	Meta     *SiteMeta       `json:"meta"`
	Elements []*SiteElements `json:"elements"`
}

func main() {
	http.HandleFunc("/", parseIt)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func parseIt(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	urls := strings.Split(string(body), "\n")
	responsesChan := make(chan *SiteData)
	finishChan := make(chan bool)
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go parser(url, &wg, responsesChan)
	}

	go func() {
		wg.Wait()
		finishChan <- true
	}()

	result := []*SiteData{}
	for {
		select {
		case res := <- responsesChan:
			result = append(result, res)
		case <- finishChan:
			j,_ := json.MarshalIndent(result, "", "  ")
			w.Write(j)
			return
		}
	}
}

func parser(url string, wgPtr *sync.WaitGroup, responsesChan chan *SiteData) {
	defer wgPtr.Done()
	url = strings.TrimSpace(url)
	if url == "" {
		return
	}
	res, err := http.Get(url)
	if err != nil {
		log.Println(url, err)
	} else {
		defer res.Body.Close()

		siteData := &SiteData{Url:url}
		headers := map[string]string{}
		for k, v := range res.Header {
			headers[strings.ToLower(k)] = string(v[0])
		}
		siteMeta := &SiteMeta{}
		siteMeta.Headers = append(siteMeta.Headers, headers)
		siteMeta.Status = res.StatusCode
		siteData.Meta = siteMeta

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(url, err)
		} else {
			tags := map[string]int{}
			re := regexp.MustCompile(`<(\w+)`)
			matches := re.FindAllStringSubmatch(string(body), -1)
			for _, sub := range matches {
				t := strings.ToLower(sub[1])
				if _, ok := tags[t]; ok {
					tags[t] += 1
				} else {
					tags[t] = 1
				}
			}
			for tagName, count := range tags {
				siteData.Elements = append(siteData.Elements, &SiteElements{TagName:tagName, Count:count})
			}
		}
		responsesChan <- siteData
	}
}

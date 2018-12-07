package monitor

import (
	"encoding/json"
	"fmt"
	"freshping/queue"
	"net/http"
	"time"
)

var purgingStaleDBDataChannel = make(chan bool)

func init() {

}

// Response is the general structure for a response
type Response struct {
	Error      *Error      `json:"Error"`
	Result     interface{} `json:"Result"`
	Success    bool        `json:"Success"`
	HTTPStatus int         `json:"HttpStatus"`
}

var urlstruct []urlstructure

//urlstructure
type urlstructure struct {
	Url       string
	Frequency int
	UrlChan   <-chan time.Time
}

var httpResponceChan chan *HTTPResponse

// HTTPResponse defines the results from a http request
type HTTPResponse struct {
	URL      string
	Response *http.Response
	Error    error
}

type URL_Monitor struct {
	Url       string
	Frequency int
	UrlChan   <-chan time.Time
}

func StartURLParser(url []URL_Monitor) {
	purgingQueueChannel = make(chan bool)

	jsonConf := util.ReadFile(configrationFile)

	if err := json.NewDecoder(strings.NewReader(jsonConf)).Decode(&urlstruct); err != nil {
		fmt.Println("Error Json Decording : ", err)
	}

	for _, url := range urlstruct {
		url.UrlChan = time.Tick(time.Duration(time.url.Frequency) * time.Minute)
	}

}

func GetURLHttpResponce(url) {

	fmt.Printf("Fetching %s \n", url)
	trans := &http.Transport{ResponseHeaderTimeout: time.Duration(5 * time.Second), DisableKeepAlives: true}
	client := &http.Client{
		Transport: trans,
		Timeout:   time.Duration(5 * time.Second),
	}

	resp, err := client.Get(url)
	defer resp.Body.Close()
	queue.Queue.Enqueue(&HTTPResponse{url, resp, err})
	purgingQueueChannel <- true
}

func pushDataElasticSearch(data interface{}) {

}

func sendDatatoElasticSearch() {

	for {
		time.Sleep(1 * time.Nanosecond)

		select {
		case <-url.UrlChan:
			go GetURLHttpResponce(url)
		case <-purgingQueueChannel:
			for queue.Queue.Len() > 0 {
				go pushDataElasticSearch(queue.Queue.Dequeue())
			}

		default:
		}
	}
}

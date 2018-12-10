package monitor

import (
	"encoding/json"
	"fmt"
	"freshping/elasticsearch"
	"freshping/queue"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var (
	configrationFile    string = "/Users/maropost/Desktop/maropost/src/freshping/config/monitorJson.json"
	purgingQueueChannel chan bool

	FreshPing        freshping
	httpResponceChan chan *HTTPResponse
)

// HTTPResponse defines the results from a http request
type HTTPResponse struct {
	ClientID  string
	URL       string
	Response  *http.Response
	Error     error
	TimeStamp int64
}

type urldata struct {
	URL       string `json:"url"`
	Frequency int    `json:"frequency"`
}

type urlmonitor struct {
	ClientID string     `json:"client_id"`
	URLData  []*urldata `json:"url_data"`
}

type freshping struct {
	URLMonitor []*urlmonitor `json:"freshping"`
}

type app_Monitor struct {
	ClientID string
	URL      string
	Appmons  <-chan time.Time
}

var ApplicationMonitor []app_Monitor

func StartURLParser() {
	purgingQueueChannel = make(chan bool)

	jsonConf, errs := ioutil.ReadFile(configrationFile)
	if errs != nil {
		fmt.Println("Error ReadFile : ", errs)
	}

	if err := json.NewDecoder(strings.NewReader(string(jsonConf[:]))).Decode(&FreshPing); err != nil {

	}

	for _, urlmonitor := range FreshPing.URLMonitor {
		for _, url := range urlmonitor.URLData {
			var appMonit app_Monitor
			appMonit.Appmons = time.Tick(time.Duration(url.Frequency) * time.Minute)
			appMonit.ClientID = urlmonitor.ClientID
			appMonit.URL = url.URL
			ApplicationMonitor = append(ApplicationMonitor, appMonit)
		}
	}

	go sendDatatoElasticSearch()
}

func getURLHttpResponce(url string, ClientID string) {

	//fmt.Printf("Fetching %s \n", url)
	trans := &http.Transport{ResponseHeaderTimeout: time.Duration(5 * time.Second), DisableKeepAlives: true}
	client := &http.Client{
		Transport: trans,
		Timeout:   time.Duration(5 * time.Second),
	}

	resp, err := client.Get("http://" + url)
	//defer resp.Body.Close()
	fmt.Println(url, "    resp : ", resp.Status, " : err : ", err)
	queue.Queue.Enqueue(&HTTPResponse{ClientID: ClientID, URL: url, Response: resp, Error: err, TimeStamp: time.Now().UnixNano()})
	purgingQueueChannel <- true
}

func pushDataElasticSearch(data interface{}) {
	value := data.(*HTTPResponse)
	fmt.Println("value : ", value)
	output, err := json.Marshal(*value)
	if err != nil {
		fmt.Println("error json.Marshal : ", err)
	}
	fmt.Println("Type : ", reflect.TypeOf(data))
	elasticsearch.InsertDataToElastic(value.ClientID, string(output[:]))
}

func sendDatatoElasticSearch() {
	fmt.Println("ApplicationMonitor : ", ApplicationMonitor)
	for {
		for _, urldata := range ApplicationMonitor {
			select {
			case <-urldata.Appmons:
				go getURLHttpResponce(urldata.URL, urldata.ClientID)
			case <-purgingQueueChannel:
				for queue.Queue.Len() > 0 {
					go pushDataElasticSearch(queue.Queue.Dequeue())
				}

			default:
			}
		}
		time.Sleep(1 * time.Nanosecond)
	}

}

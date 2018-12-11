package monitor

import (
	"encoding/json"
	"fmt"
	"freshping/elasticsearch"
	"freshping/queue"
	"freshping/util"
	"io/ioutil"
	"net/http"
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
	ClientID string
	URL      string
	//Response  *http.Response
	Response  int
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
	go QueueStatus()
}

func QueueStatus() {
	ch := time.Tick(time.Duration(30) * time.Second)

	for {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-ch:
			if queue.Queue.Len() > 0 {
				purgingQueueChannel <- true
			}
		}
	}
}

func getURLHttpResponce(url string, ClientID string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover getURLHttpResponce: ", util.RecoverExceptionDetails(util.FuncName()), "  : ", err)
		}
	}()
	//fmt.Printf("Fetching %s \n", url)
	trans := &http.Transport{ResponseHeaderTimeout: time.Duration(5 * time.Second), DisableKeepAlives: true}
	client := &http.Client{
		Transport: trans,
		Timeout:   time.Duration(5 * time.Second),
	}

	resp, err := client.Get("http://" + url)
	//defer resp.Body.Close()
	fmt.Println(url, "    resp : ", resp, " : err : ", err)
	httpstatus := 0
	if resp != nil {
		httpstatus = resp.StatusCode
	}
	queue.Queue.Enqueue(&HTTPResponse{ClientID: ClientID, URL: url, Error: err, Response: httpstatus, TimeStamp: time.Now().UnixNano()})

}

func pushDataElasticSearch() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover pushDataElasticSearch: ", util.RecoverExceptionDetails(util.FuncName()), "   ", err)
		}
	}()

	for queue.Queue.Len() > 0 {
		data := queue.Queue.Dequeue()
		value := data.(*HTTPResponse)

		//fmt.Println("value : ", value)
		output, err := json.Marshal(data)
		if err != nil {
			fmt.Println("error json.Marshal : ", err)
		}
		strtime := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
		fmt.Println(strtime, "  | output : ", string(output[:]))
		elasticsearch.InsertDataToElastic(value.ClientID, string(output[:]))
	}
	//purgingQueueChannel <- false
}

func sendDatatoElasticSearch() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover sendDatatoElasticSearch: ", util.RecoverExceptionDetails(util.FuncName()), "    ", err)
		}
	}()
	fmt.Println("ApplicationMonitor : ", ApplicationMonitor)

	for {
		for _, urldata := range ApplicationMonitor {
			select {
			case <-urldata.Appmons:
				go getURLHttpResponce(urldata.URL, urldata.ClientID)
			case <-purgingQueueChannel:
				go pushDataElasticSearch()
				purgingQueueChannel <- false
			}
			time.Sleep(1 * time.Nanosecond)
		}
		time.Sleep(1 * time.Nanosecond)
	}
}

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"
)

var optThreads = flag.Int("threads", runtime.NumCPU(), "Number of system threads")
var optWorkers = flag.Int("workers", 20, "Number of HTTP GET Workers")
var optUrl = flag.String("url", "http://127.0.0.1:8080/", "URL for HTTP GET Method")

type UrlRequest struct {
	url  string
	pipe chan *[]byte
}

type HttpClient struct {
	request chan *UrlRequest
}

func NewHttpClient(workers int) *HttpClient {
	request := make(chan *UrlRequest, 10*workers)

	tp := &http.Transport{
		MaxIdleConnsPerHost:   workers,
		ResponseHeaderTimeout: 15 * time.Second,
	}
	client := &http.Client{Transport: tp}

	worker := func() {
		for {
			req := <-request
			res, err := client.Get(req.url)
			if err != nil {
				log.Println(err)
				req.pipe <- nil
				continue
			}

			raw, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Println(err)
				req.pipe <- nil
				continue
			}

			req.pipe <- &raw
		}
	}

	for i := 0; i < workers; i++ {
		go worker()
	}

	return &HttpClient{
		request: request,
	}
}

func (this *HttpClient) Get(url string) *[]byte {
	reqRaw := &UrlRequest{url, make(chan *[]byte, 1)}
	this.request <- reqRaw
	return <-reqRaw.pipe
}

func main() {
	fmt.Println("HTTP Client")

	flag.Parse()

	fmt.Println("Threads:", *optThreads)
	runtime.GOMAXPROCS(*optThreads)

	fmt.Println("Workers:", *optWorkers)
	httpClient := NewHttpClient(*optWorkers)

	fmt.Println("HTTP GET:", *optUrl)
	raw := httpClient.Get(*optUrl)
	fmt.Println(string(*raw))
}

package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"
)

var optThreads = flag.Int("threads", runtime.NumCPU(), "Max system threads (default number of CPUs)")
var optBindAddress = flag.String("bind", "127.0.0.1:8080", "HTTP Server bind address (default 127.0.0.1:8080)")

type Parameters map[string][]string

func ParseRequestParameters(req *http.Request) *Parameters {
	params := Parameters(req.URL.Query())
	return &params
}

type HttpHandler struct{}

func (*HttpHandler) ServeHTTP(out http.ResponseWriter, req *http.Request) {
	ParseRequestParameters(req)
	out.Write([]byte("That is it!\n"))
}

func main() {
	log.Println("HTTP Server")

	flag.Parse()

	log.Println("Threads:", *optThreads)
	runtime.GOMAXPROCS(*optThreads)

	log.Println("Address:", *optBindAddress)
	err := http.ListenAndServe(*optBindAddress, &HttpHandler{})
	if err != nil {
		log.Fatal("Error starting HTTP Server", err)
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

type Configuration struct {
	HttpBindAddress string `json:"http_bindaddress"`
}

func (this *Configuration) String() string {
	var out bytes.Buffer

	fmt.Fprintln(&out, "HttpBindAddress:", this.HttpBindAddress)

	return out.String()
}

func LoadConfiguration(file string) *Configuration {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("Error reading configuration file", file, err)
		return nil
	}
	var c Configuration
	json.Unmarshal(raw, &c)
	return &c
}

var config = LoadConfiguration("conf/config.json")

type HttpHandler struct{}

func (this HttpHandler) ServeHTTP(out http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(out, "May I help you?")
}

func main() {
	fmt.Println("Espresso!")

	ncpu := runtime.NumCPU()
	runtime.GOMAXPROCS(ncpu)
	fmt.Println("CPUs:", ncpu)

	fmt.Println("(Configuration)")
	fmt.Println(config)

	fmt.Println("HTTP Server start!")
	err := http.ListenAndServe(config.HttpBindAddress, &HttpHandler{})
	if err != nil {
		log.Println("Error starting HTTP Server", config.HttpBindAddress, err)
	}
}

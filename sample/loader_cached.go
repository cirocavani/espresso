package main

import (
	"container/list"
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

var optThreads = flag.Int("threads", runtime.NumCPU(), "Number of system threads")

type Data struct {
	value int
}

func (this *Data) String() string {
	out := "{"
	out += fmt.Sprintf("value: %d", this.value)
	out += "}"
	return out
}

type DataContainer []*Data

func (this *DataContainer) String() string {
	out := "["
	for i, v := range *this {
		if i != 0 {
			out += ","
		}
		out += v.String()
	}
	out += "]"
	return out
}

type Cache struct {
	sync.RWMutex
	data map[string]*DataContainer
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]*DataContainer),
	}
}

func (this *Cache) Get(key string) *DataContainer {
	this.RLock()
	defer this.RUnlock()
	value, ok := this.data[key]
	if !ok {
		log.Println("Not found", key)
		return nil
	}
	return value
}

func (this *Cache) Set(key string, value *DataContainer) *DataContainer {
	this.Lock()
	defer this.Unlock()
	this.data[key] = value
	return value
}

type DataPipe chan *DataContainer

type DataRequest struct {
	key   string
	reply DataPipe
}

type DataResult struct {
	key   string
	value *DataContainer
}

type DataLoader struct {
	cache   *Cache
	request chan *DataRequest
}

func NewDataLoader(load func(string) *DataContainer) *DataLoader {
	cache := NewCache()
	request := make(chan *DataRequest, 100)
	result := make(chan *DataResult, 100)
	work := make(map[string]*list.List)

	worker := func(key string) {
		value := load(key)
		cache.Set(key, value)
		result <- &DataResult{key, value}
	}

	requestHandler := func(req *DataRequest) {
		if pipe, ok := work[req.key]; ok {
			pipe.PushBack(&req.reply)
			return
		}
		pipe := list.New()
		pipe.PushBack(&req.reply)
		work[req.key] = pipe
		go worker(req.key)
	}

	resultHandler := func(res *DataResult) {
		pipe := work[res.key]
		for i := pipe.Front(); i != nil; i = i.Next() {
			*i.Value.(*DataPipe) <- res.value
		}
		delete(work, res.key)
	}

	broker := func() {
		for {
			select {
			case req := <-request:
				requestHandler(req)
			case res := <-result:
				resultHandler(res)
			}
		}
	}

	go broker()

	return &DataLoader{
		cache:   cache,
		request: request,
	}
}

func (this *DataLoader) load(key string) DataPipe {
	pipe := make(DataPipe, 1)
	this.request <- &DataRequest{key, pipe}
	return pipe
}

func (this *DataLoader) fetch(key string, timeout time.Duration) *DataContainer {
	pipe := this.load(key)

	if timeout == 0 {
		return <-pipe
	}

	select {
	case value := <-pipe:
		return value
	case <-time.After(timeout):
		return nil
	}
}

func (this *DataLoader) get(key string) *DataContainer {
	return this.cache.Get(key)
}

func (this *DataLoader) Fetch(key string, timeout time.Duration) *DataContainer {
	if value := this.get(key); value != nil {
		return value
	}
	return this.fetch(key, timeout)
}

func main() {
	fmt.Println("Loader with Cache")

	flag.Parse()

	fmt.Println("Threads:", *optThreads)
	runtime.GOMAXPROCS(*optThreads)

	loader := NewDataLoader(func(key string) *DataContainer {
		log.Println("Loading:", key)
		time.Sleep(5 * time.Second)
		return &DataContainer{
			&Data{1},
			&Data{2},
		}
	})

	timeout := 1 * time.Second

	v := loader.Fetch("x", timeout)
	log.Println("Fetch 1:", v)

	v = loader.Fetch("x", timeout)
	log.Println("Fetch 2:", v)

	v = loader.Fetch("x", 0)
	log.Println("Fetch 3:", v)

	v = loader.Fetch("x", timeout)
	log.Println("Fetch 4:", v)

	v = loader.Fetch("x", timeout)
	log.Println("Fetch 5:", v)
}

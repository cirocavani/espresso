package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

var optThreads = flag.Int("threads", runtime.NumCPU(), "Number of system threads")

type EventType int

const (
	SET EventType = iota
	DEL
)

type Event struct {
	Type      EventType
	Key       string
	Timestamp time.Time
}

func (this *Event) String() string {
	return fmt.Sprintf("%v, %s, %v", this.Type, this.Key, this.Timestamp)
}

type Cache struct {
	sync.RWMutex
	data map[string]interface{}

	event chan *Event
}

func NewCache(ttl time.Duration, size int) *Cache {
	event := make(chan *Event, 100)

	eviction := func() {
		for {
			e := <-event
			fmt.Println(e)
		}
	}

	go eviction()

	return &Cache{
		data:  make(map[string]interface{}),
		event: event,
	}
}

func (this *Cache) String() string {
	this.RLock()
	defer this.RUnlock()

	out := "{"
	first := true
	for k, v := range this.data {
		if first {
			first = false
		} else {
			out += ","
		}
		out += fmt.Sprintf(`"%s"=%#v`, k, v)
	}
	out += "}"

	return out
}

func (this *Cache) Get(key string) interface{} {
	this.RLock()
	defer this.RUnlock()
	value, ok := this.data[key]
	if !ok {
		log.Println("Not found", key)
		return nil
	}
	return value
}

func (this *Cache) riseEvent(t EventType, keys ...string) {
	timestamp := time.Now()
	go func() {
		for _, key := range keys {
			this.event <- &Event{t, key, timestamp}
		}
	}()
}

func (this *Cache) Set(key string, value interface{}) interface{} {
	this.Lock()
	defer this.Unlock()
	defer this.riseEvent(SET, key)
	this.data[key] = value
	return value
}

func (this *Cache) Delete(keys ...string) {
	this.Lock()
	defer this.Unlock()
	defer this.riseEvent(DEL, keys...)
	for _, key := range keys {
		delete(this.data, key)
	}
}

func main() {
	fmt.Println("Cache")

	flag.Parse()

	fmt.Println("Threads:", *optThreads)
	runtime.GOMAXPROCS(*optThreads)

	ttl := 5 * time.Second
	size := 10

	cache := NewCache(ttl, size)
	cache.Set("x", []int{1, 2, 3})
	cache.Set("y", "[1 2 3]")

	fmt.Println(cache)

	cache.Delete("y")

	fmt.Println(cache)

	time.Sleep(10 * time.Second)
}

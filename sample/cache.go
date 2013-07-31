package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync"
	_ "time"
)

var optThreads = flag.Int("threads", runtime.NumCPU(), "Number of system threads")

type Cache struct {
	sync.RWMutex
	data map[string]interface{}
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
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

func (this *Cache) Set(key string, value interface{}) interface{} {
	this.Lock()
	defer this.Unlock()
	this.data[key] = value
	return value
}

func (this *Cache) Delete(keys ...string) {
	this.Lock()
	defer this.Unlock()
	for _, key := range keys {
		delete(this.data, key)
	}
}

func main() {
	fmt.Println("Cache")

	flag.Parse()

	fmt.Println("Threads:", *optThreads)
	runtime.GOMAXPROCS(*optThreads)

	cache := NewCache()
	cache.Set("x", []int{1, 2, 3})
	cache.Set("y", "[1 2 3]")

	fmt.Println(cache)

	cache.Delete("y")

	fmt.Println(cache)
}

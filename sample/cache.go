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

type CacheValue struct {
	Value      interface{}
	expiration time.Time
}

type Cache struct {
	sync.RWMutex
	data    map[string]*CacheValue
	ttl     time.Duration
	maxSize int
}

func NewCache(ttl time.Duration, maxSize int) *Cache {
	cache := &Cache{
		data:    make(map[string]*CacheValue),
		ttl:     ttl,
		maxSize: maxSize,
	}

	eviction := func() {
		ticker := time.NewTicker(ttl)

		for {
			<-ticker.C
			cache.Eviction()
		}
	}

	go eviction()

	return cache
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

func (this *Cache) keysOlderThan(t *time.Time) *list.List {
	this.RLock()
	defer this.RUnlock()

	keys := list.New()
	for k, v := range this.data {
		if v.expiration.Before(*t) {
			keys.PushBack(k)
		}
	}
	return keys
}

func (this *Cache) deleteKeys(keys *list.List) {
	this.Lock()
	defer this.Unlock()

	for i := keys.Front(); i != nil; i = i.Next() {
		delete(this.data, i.Value.(string))
	}
}

func (this *Cache) removeExpired() {
	now := time.Now()
	keys := this.keysOlderThan(&now)
	if keys.Len() > 0 {
		go this.deleteKeys(keys)
	}
}

func (this *Cache) Eviction() {
	this.removeExpired()
}

func (this *Cache) Get(key string) interface{} {
	this.RLock()
	defer this.RUnlock()

	value, ok := this.data[key]
	if !ok {
		log.Println("Not found", key)
		return nil
	}
	return value.Value
}

func (this *Cache) Set(key string, value interface{}) interface{} {
	this.Lock()
	defer this.Unlock()

	this.data[key] = &CacheValue{value, time.Now().Add(this.ttl)}
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

	ttl := 5 * time.Second
	size := 10

	cache := NewCache(ttl, size)
	cache.Set("x", []int{1, 2, 3})
	cache.Set("y", "[1 2 3]")

	fmt.Println("Cache TTL 5s: with x, y")
	fmt.Println(cache)

	cache.Delete("y")

	fmt.Println("Cache TTL 5s: y deleted")
	fmt.Println(cache)

	time.Sleep(10 * time.Second)

	fmt.Println("Cache TTL 5s: after 10s")
	fmt.Println(cache)
}

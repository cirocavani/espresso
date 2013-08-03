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

type CacheEntry struct {
	Key        *string
	Value      interface{}
	expiration *time.Time
}

func (this *CacheEntry) Before(t *time.Time) bool {
	return this.expiration.Before(*t)
}

type Cache struct {
	sync.RWMutex
	data map[string]*CacheEntry
	ttl  time.Duration
}

func NewCache(purge, ttl time.Duration) *Cache {
	cache := &Cache{
		data: make(map[string]*CacheEntry),
		ttl:  ttl,
	}

	eviction := func() {
		ticker := time.NewTicker(purge)

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

func (this *Cache) entriesOlderThan(t *time.Time) *list.List {
	this.RLock()
	defer this.RUnlock()

	entries := list.New()
	for _, v := range this.data {
		if v.Before(t) {
			entries.PushBack(v)
		}
	}
	return entries
}

func (this *Cache) deleteEntries(entries *list.List) {
	this.Lock()
	defer this.Unlock()

	for i := entries.Front(); i != nil; i = i.Next() {
		vi := i.Value.(*CacheEntry)
		if v, ok := this.data[*vi.Key]; ok && vi == v {
			delete(this.data, *vi.Key)
		}
	}
}

func (this *Cache) removeExpired() {
	now := time.Now()
	entries := this.entriesOlderThan(&now)
	if entries.Len() > 0 {
		go this.deleteEntries(entries)
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

func (this *Cache) SetExpiration(key string, value interface{}, ttl time.Duration) {
	this.Lock()
	defer this.Unlock()

	expiration := time.Now().Add(ttl)
	this.data[key] = &CacheEntry{
		&key,
		value,
		&expiration,
	}
}

func (this *Cache) Set(key string, value interface{}) {
	this.SetExpiration(key, value, this.ttl)
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

	purge := 1 * time.Second
	ttl := 5 * time.Second

	cache := NewCache(purge, ttl)
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

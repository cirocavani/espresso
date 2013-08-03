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
	data    map[string]*list.Element
	entries *list.List

	ttl     time.Duration
	maxSize int
}

func NewCache(purge, ttl time.Duration, maxSize int) *Cache {
	cache := &Cache{
		data:    make(map[string]*list.Element),
		entries: list.New(),
		ttl:     ttl,
		maxSize: maxSize,
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

func (this *Cache) Size() int {
	this.RLock()
	defer this.RUnlock()

	return this.entries.Len()
}

func (this *Cache) removeEntry(i *list.Element) {
	value := i.Value.(*CacheEntry)
	delete(this.data, *value.Key)
	this.entries.Remove(i)
}

func (this *Cache) removeExpired() {
	now := time.Now()
	for i := this.entries.Front(); i != nil; {
		next := i.Next()
		if i.Value.(*CacheEntry).Before(&now) {
			this.removeEntry(i)
		}
		i = next
	}
}

func (this *Cache) removeOverflow() {
	for over := this.entries.Len() - this.maxSize; over > 0; over-- {
		i := this.entries.Front()
		this.removeEntry(i)
	}
}

func (this *Cache) Eviction() {
	this.Lock()
	defer this.Unlock()

	this.removeExpired()
	this.removeOverflow()
}

func (this *Cache) Get(key string) interface{} {
	this.RLock()
	defer this.RUnlock()

	i, ok := this.data[key]
	if !ok {
		log.Println("Not found", key)
		return nil
	}
	value := i.Value.(*CacheEntry)
	return value.Value
}

func (this *Cache) release(key string) {
	if i, ok := this.data[key]; ok {
		delete(this.data, key)
		this.entries.Remove(i)
	}
}

func (this *Cache) Set(key string, value interface{}) {
	this.SetExpiration(key, value, this.ttl)
}

func (this *Cache) SetExpiration(key string, value interface{}, ttl time.Duration) {
	this.Lock()
	defer this.Unlock()

	this.release(key)

	expiration := time.Now().Add(ttl)
	entry := &CacheEntry{
		&key,
		value,
		&expiration,
	}
	this.data[key] = this.entries.PushBack(entry)
}

func (this *Cache) Release(keys ...string) {
	this.Lock()
	defer this.Unlock()

	for _, key := range keys {
		this.release(key)
	}
}

func main() {
	fmt.Println("Cache")

	flag.Parse()

	fmt.Println("Threads:", *optThreads)
	runtime.GOMAXPROCS(*optThreads)

	purge := 1 * time.Second
	ttl := 5 * time.Second
	size := 10

	cache := NewCache(purge, ttl, size)
	cache.Set("x", []int{1, 2, 3})
	cache.Set("y", "[1 2 3]")

	fmt.Println("Cache TTL 5s: with x, y")
	fmt.Println(cache)

	cache.Release("y")

	fmt.Println("Cache TTL 5s: y deleted")
	fmt.Println(cache)

	time.Sleep(10 * time.Second)

	fmt.Println("Cache TTL 5s: after 10s")
	fmt.Println(cache)

	v := "some very complex data"
	for _, k := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"} {
		cache.SetExpiration(k, &v, 10*time.Hour)
	}

	fmt.Println("Cache Size:", cache.Size())

	time.Sleep(10 * time.Second)

	fmt.Println("Cache Size after 10s:", cache.Size())
}

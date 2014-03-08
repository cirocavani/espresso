package main

import (
	"container/list"
	"flag"
	"fmt"
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
	entries *list.List
	index   map[string]*list.Element

	ttl     time.Duration
	maxSize int
}

func NewCache(purge, ttl time.Duration, maxSize int) *Cache {
	cache := &Cache{
		entries: list.New(),
		index:   make(map[string]*list.Element),
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
	for i := this.entries.Front(); i != nil; i = i.Next() {
		if first {
			first = false
		} else {
			out += ","
		}
		entry := i.Value.(*CacheEntry)
		out += fmt.Sprintf(`"%s"=%#v`, *entry.Key, entry.Value)
	}
	out += "}"

	return out
}

func (this *Cache) size() int {
	return this.entries.Len()
}

func (this *Cache) Size() int {
	this.RLock()
	defer this.RUnlock()

	return this.size()
}

func (this *Cache) removeEntry(i *list.Element) {
	entry := i.Value.(*CacheEntry)
	this.entries.Remove(i)
	delete(this.index, *entry.Key)
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
	for over := this.size() - this.maxSize; over > 0; over-- {
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

func (this *Cache) Get(key string) (interface{}, bool) {
	this.RLock()
	defer this.RUnlock()

	i, ok := this.index[key]
	if !ok {
		return nil, false
	}
	entry := i.Value.(*CacheEntry)
	return entry.Value, true
}

func (this *Cache) release(key string) {
	if i, ok := this.index[key]; ok {
		this.entries.Remove(i)
		delete(this.index, key)
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
	this.index[key] = this.entries.PushBack(entry)
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

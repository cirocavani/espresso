package main

import (
	"flag"
	"fmt"
	"runtime"
)

import "github.com/hoisie/redis"

var optThreads = flag.Int("threads", runtime.NumCPU(), "Number of system threads")
var optServer = flag.String("server", "127.0.0.1:6379", "Redis server address (ip:port)")
var optPoolSize = flag.Int("pool", 20, "Redis connection pool size")

type RedisDriver struct {
	client redis.Client
}

func NewRedisDriver(address string, poolSize int) *RedisDriver {
	client := redis.Client{
		Addr:        address,
		MaxPoolSize: poolSize,
	}
	return &RedisDriver{
		client: client,
	}
}

func (this *RedisDriver) Get(key string) (*[]byte, error) {
	raw, err := this.client.Get(key)
	return &raw, err
}

func (this *RedisDriver) Set(key string, value *[]byte) {
	this.client.Set(key, *value)
}

func main() {
	fmt.Println("Redis")

	flag.Parse()

	fmt.Println("Threads:", *optThreads)
	runtime.GOMAXPROCS(*optThreads)

	fmt.Println("Redis:", *optServer)
	fmt.Println("Pool:", *optPoolSize)

	client := NewRedisDriver(*optServer, *optPoolSize)

	value := []byte("This is a very data-intensive computing value result")
	client.Set("x", &value)
	raw, _ := client.Get("x")
	fmt.Println(string(*raw))
}

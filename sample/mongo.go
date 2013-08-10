package main

import (
	"flag"
	"fmt"
	"runtime"
)

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var optThreads = flag.Int("threads", runtime.NumCPU(), "Number of system threads")
var optServer = flag.String("server", "127.0.0.1:27017", "MongoDB Server address")
var optDatabase = flag.String("database", "test", "MongoDB Database Name")
var optCollection = flag.String("collection", "test", "MongoDB Collection Name")

type MongoDriver struct {
	session    *mgo.Session
	database   string
	collection string
}

func NewMongoDriver(address, database, collection string) *MongoDriver {
	session, err := mgo.Dial(address)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return &MongoDriver{session, database, collection}
}

func (this *MongoDriver) Close() {
	if this.session == nil {
		return
	}
	this.session.Close()
	this.session = nil
}

func (this *MongoDriver) c() *mgo.Collection {
	if this.session == nil {
		return nil
	}
	db := this.session.DB(this.database)
	return db.C(this.collection)
}

func (this *MongoDriver) Store(values ...interface{}) {
	c := this.c()
	err := c.Insert(values...)
	if err != nil {
		panic(err)
	}
}

func (this *MongoDriver) Fetch(query bson.M, result interface{}) {
	c := this.c()
	err := c.Find(query).One(result)
	if err != nil {
		panic(err)
	}
}

func (this *MongoDriver) Delete(id interface{}) {
	c := this.c()
	err := c.RemoveId(id)
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Mongo")

	flag.Parse()

	fmt.Println("Threads:", *optThreads)
	runtime.GOMAXPROCS(*optThreads)

	fmt.Println("MongoDB:", *optServer)
	fmt.Println("Database:", *optDatabase)
	fmt.Println("Collection:", *optCollection)
	mongo := NewMongoDriver(*optServer, *optDatabase, *optCollection)

	fmt.Println("Data to MongoDB")
	value := map[string]interface{}{"a": "b", "c": 10}
	fmt.Printf("%+v\n", value)
	mongo.Store(&value)

	fmt.Println("Data from MongoDB")
	result := make(map[string]interface{})
	mongo.Fetch(bson.M{"a": "b"}, &result)
	fmt.Printf("%+v\n", result)

	fmt.Printf("Removing '%v'\n", result["_id"])
	mongo.Delete(result["_id"])
}

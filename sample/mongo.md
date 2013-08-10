MongoDB Server:

	wget http://fastdl.mongodb.org/linux/mongodb-linux-x86_64-2.4.5.tgz
	tar xzf mongodb-linux-x86_64-2.4.5.tgz
	cd mongodb-linux-x86_64-2.4.5
	mkdir data
	bin/mongod --dbpath data

MongoDB Client:

http://labix.org/mgo

http://godoc.org/labix.org/v2/mgo
http://godoc.org/labix.org/v2/mgo/bson

	(Bzr is required)
	export GOPATH=<SOURCE FOLDER>/espresso/lib
	go get labix.org/v2/mgo

Output:

	go run sample/mongo.go
	
	Mongo
	Threads: 8
	MongoDB: 127.0.0.1:27017
	Database: test
	Collection: test
	Data to MongoDB
	map[a:b c:10]
	Data from MongoDB
	map[_id:ObjectIdHex("5205b56037278f48a8519e44") a:b c:10]
	Removing 'ObjectIdHex("5205b56037278f48a8519e44")'

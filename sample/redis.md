Redis Server:

      wget http://redis.googlecode.com/files/redis-2.6.14.tar.gz
      tar xzf redis-2.6.14.tar.gz
      cd redis-2.6.14
      make
      
      src/redis-server

Redis Client:

      export GOPATH=<SOURCE FOLDER>/espresso/lib
      go get github.com/hoisie/redis

Output:

	go run sample/redis.go
	
	Redis
	Threads: 8
	Redis: 127.0.0.1:6379
	Pool: 20
	This is a very data-intensive computing value result

Server:

	go run src/sample/http_request/main.go 
	2013/07/25 00:21:28 HTTP Server
	2013/07/25 00:21:28 Threads: 8
	2013/07/25 00:21:28 Address: 127.0.0.1:8080

Benchmark:

	wrk -t 15 -c 600 -d 20s "http://127.0.0.1:8080/?a=1&b=2&a=3&c=4&c=5"
	Running 20s test @ http://127.0.0.1:8080/?a=1&b=2&a=3&c=4&c=5
	  15 threads and 600 connections
	  Thread Stats   Avg      Stdev     Max   +/- Stdev
	    Latency     7.14ms    2.99ms 219.16ms   78.15%
	    Req/Sec     5.77k     1.33k   17.64k    70.43%
	  1652237 requests in 20.00s, 203.26MB read
	Requests/sec:  82632.14
	Transfer/sec:     10.17MB

	wrk -t 15 -c 600 -d 20s "http://127.0.0.1:8080/?a=1&b=2&a=3&c=4&c=5"
	Running 20s test @ http://127.0.0.1:8080/?a=1&b=2&a=3&c=4&c=5
	  15 threads and 600 connections
	  Thread Stats   Avg      Stdev     Max   +/- Stdev
	    Latency     7.09ms    3.17ms 226.98ms   81.09%
	    Req/Sec     5.74k     1.12k   12.54k    70.55%
	  1665390 requests in 20.00s, 204.88MB read
	Requests/sec:  83269.48
	Transfer/sec:     10.24MB

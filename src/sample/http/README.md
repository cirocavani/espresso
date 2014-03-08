Server:

	go run src/sample/http/main.go
	2013/07/24 21:00:22 HTTP Server
	2013/07/24 21:00:22 Threads: 8
	2013/07/24 21:00:22 Address: 127.0.0.1:8080

Benchmark:

	wrk -t 15 -c 600 -d 20s http://127.0.0.1:8080/
	Running 20s test @ http://127.0.0.1:8080/
	  15 threads and 600 connections
	  Thread Stats   Avg      Stdev     Max   +/- Stdev
	    Latency     6.47ms    4.00ms 212.88ms   91.46%
	    Req/Sec     6.42k     1.53k   23.00k    70.53%
	  1829292 requests in 20.00s, 225.05MB read
	Requests/sec:  91484.37
	Transfer/sec:     11.25MB

	wrk -t 15 -c 600 -d 20s http://127.0.0.1:8080/
	Running 20s test @ http://127.0.0.1:8080/
	  15 threads and 600 connections
	  Thread Stats   Avg      Stdev     Max   +/- Stdev
	    Latency     6.45ms    2.44ms  35.74ms   74.71%
	    Req/Sec     6.37k     1.33k   17.40k    70.70%
	  1828908 requests in 20.00s, 225.00MB read
	Requests/sec:  91454.75
	Transfer/sec:     11.25MB

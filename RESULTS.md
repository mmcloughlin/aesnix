```
go version go1.8.1 darwin/amd64
=== RUN   TestSingle
--- PASS: TestSingle (0.00s)
=== RUN   TestMulti
--- PASS: TestMulti (0.00s)
BenchmarkSingle-4   	1000000000	         6.84 ns/op	2340.60 MB/s
BenchmarkMulti/2-4  	1000000000	         8.57 ns/op	3732.92 MB/s
BenchmarkMulti/4-4  	500000000	        12.6 ns/op	5077.95 MB/s
BenchmarkMulti/6-4  	500000000	        18.8 ns/op	5113.87 MB/s
BenchmarkMulti/8-4  	300000000	        24.3 ns/op	5265.70 MB/s
BenchmarkMulti/10-4 	200000000	        31.2 ns/op	5131.06 MB/s
BenchmarkMulti/12-4 	200000000	        37.8 ns/op	5084.23 MB/s
BenchmarkMulti/14-4 	200000000	        43.4 ns/op	5161.76 MB/s
PASS
```

# aesnix

Experiments with AES-NI performance in Golang.

These tests support the implementation of [AES-CTR
mode](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation) in
assembly (Golang issue [#20967](https://golang.org/issue/20967)). We are
interested in performance improvements obtained by encrypting multiple blocks
at the same time (with the same key).

Results show massive improvements up to 4 concurrent blocks and inconsistent
or diminishing returns beyond that. Comparisons to functions that purely
contain `AESENC` calls suggest the other memory and arithmetic operations are
lost in the noise (as you would expect).

```
$ go test -bench . -benchtime 20s
BenchmarkSingle-4   	5000000000	         7.08 ns/op	2258.35 MB/s
BenchmarkMulti/2-4  	5000000000	         7.75 ns/op	4128.89 MB/s
BenchmarkMulti/4-4  	2000000000	        12.8 ns/op	5007.81 MB/s
BenchmarkMulti/6-4  	2000000000	        19.2 ns/op	5010.51 MB/s
BenchmarkMulti/8-4  	1000000000	        25.7 ns/op	4975.29 MB/s
BenchmarkMulti/10-4 	1000000000	        32.1 ns/op	4987.93 MB/s
BenchmarkMulti/12-4 	1000000000	        38.4 ns/op	5000.29 MB/s
BenchmarkMulti/14-4 	1000000000	        44.6 ns/op	5025.15 MB/s
BenchmarkNomem/2-4  	5000000000	         6.38 ns/op	5018.09 MB/s
BenchmarkNomem/4-4  	2000000000	        12.8 ns/op	5000.03 MB/s
BenchmarkNomem/6-4  	2000000000	        19.1 ns/op	5023.91 MB/s
BenchmarkNomem/8-4  	1000000000	        25.8 ns/op	4963.32 MB/s
BenchmarkNomem/10-4 	1000000000	        32.0 ns/op	5005.75 MB/s
BenchmarkNomem/12-4 	1000000000	        38.6 ns/op	4978.14 MB/s
BenchmarkNomem/14-4 	1000000000	        44.9 ns/op	4989.37 MB/s
```

See [results of some more in-depth tests](RESULTS.md).

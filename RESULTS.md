## No Memory

Assembly functions that make the same number of `AESENC` calls but do not
access memory. Initially confusing results.

```
$ go version
go version go1.8.1 darwin/amd64
$ git rev-parse HEAD
09d6f625f09628977c1cb39c8a96f397b32ea308
$ go test -v -bench .
=== RUN   TestSingle
--- PASS: TestSingle (0.00s)
=== RUN   TestMulti
--- PASS: TestMulti (0.00s)
BenchmarkSingle-4   	200000000	         7.06 ns/op	2267.41 MB/s
BenchmarkMulti/2-4  	200000000	         7.61 ns/op	4202.42 MB/s
BenchmarkMulti/4-4  	100000000	        12.8 ns/op	5007.30 MB/s
BenchmarkMulti/6-4  	100000000	        19.5 ns/op	4918.08 MB/s
BenchmarkMulti/8-4  	50000000	        25.6 ns/op	4993.60 MB/s
BenchmarkMulti/10-4 	50000000	        32.6 ns/op	4902.34 MB/s
BenchmarkMulti/12-4 	50000000	        38.6 ns/op	4976.10 MB/s
BenchmarkMulti/14-4 	30000000	        45.6 ns/op	4911.10 MB/s
BenchmarkNomem/2-4  	50000000	        25.5 ns/op	1257.16 MB/s
BenchmarkNomem/4-4  	30000000	        51.6 ns/op	1240.00 MB/s
BenchmarkNomem/6-4  	20000000	        76.4 ns/op	1256.33 MB/s
BenchmarkNomem/8-4  	20000000	       103 ns/op	1242.68 MB/s
BenchmarkNomem/10-4 	10000000	       128 ns/op	1247.53 MB/s
BenchmarkNomem/12-4 	10000000	       154 ns/op	1242.72 MB/s
BenchmarkNomem/14-4 	10000000	       177 ns/op	1259.87 MB/s
PASS
```

## No Memory with Random Registers

Suspected the problem was that we used `X0` for all inputs. If we randomize
registers we get performace similar to the multi benchmark.

```
$ go version
go version go1.8.1 darwin/amd64
$ git log --oneline -n 1
35daa22 randomize registers
$ go test -v -bench .
=== RUN   TestSingle
--- PASS: TestSingle (0.00s)
=== RUN   TestMulti
--- PASS: TestMulti (0.00s)
BenchmarkSingle-4   	200000000	         7.09 ns/op	2256.68 MB/s
BenchmarkMulti/2-4  	200000000	         7.59 ns/op	4214.93 MB/s
BenchmarkMulti/4-4  	100000000	        12.8 ns/op	5009.94 MB/s
BenchmarkMulti/6-4  	100000000	        19.2 ns/op	5010.34 MB/s
BenchmarkMulti/8-4  	50000000	        25.3 ns/op	5066.54 MB/s
BenchmarkMulti/10-4 	50000000	        32.1 ns/op	4986.82 MB/s
BenchmarkMulti/12-4 	50000000	        38.2 ns/op	5020.26 MB/s
BenchmarkMulti/14-4 	30000000	        44.4 ns/op	5042.50 MB/s
BenchmarkNomem/2-4  	200000000	         6.63 ns/op	4829.12 MB/s
BenchmarkNomem/4-4  	100000000	        13.2 ns/op	4856.15 MB/s
BenchmarkNomem/6-4  	100000000	        19.2 ns/op	5003.39 MB/s
BenchmarkNomem/8-4  	50000000	        25.7 ns/op	4981.19 MB/s
BenchmarkNomem/10-4 	50000000	        32.0 ns/op	4995.10 MB/s
BenchmarkNomem/12-4 	50000000	        38.5 ns/op	4986.53 MB/s
BenchmarkNomem/14-4 	30000000	        44.6 ns/op	5024.69 MB/s
PASS
```

## No Memory Constant Key Register

If keeping registers the same is a problem, what exactly is the issue. If the
key register is the same and the block register changes, that does not cause
problems.

```
$ go version
go version go1.8.1 darwin/amd64
$ git log --oneline -n 1
341e0e5 consant key register
$ go test -v -bench .
=== RUN   TestSingle
--- PASS: TestSingle (0.00s)
=== RUN   TestMulti
--- PASS: TestMulti (0.00s)
BenchmarkSingle-4   	200000000	         6.91 ns/op	2315.51 MB/s
BenchmarkMulti/2-4  	200000000	         8.76 ns/op	3653.13 MB/s
BenchmarkMulti/4-4  	100000000	        12.9 ns/op	4961.21 MB/s
BenchmarkMulti/6-4  	100000000	        19.4 ns/op	4955.73 MB/s
BenchmarkMulti/8-4  	50000000	        25.6 ns/op	5002.28 MB/s
BenchmarkMulti/10-4 	50000000	        32.0 ns/op	4996.23 MB/s
BenchmarkMulti/12-4 	50000000	        38.0 ns/op	5049.40 MB/s
BenchmarkMulti/14-4 	30000000	        44.1 ns/op	5080.01 MB/s
BenchmarkNomem/2-4  	200000000	         6.35 ns/op	5042.83 MB/s
BenchmarkNomem/4-4  	100000000	        12.8 ns/op	5010.62 MB/s
BenchmarkNomem/6-4  	100000000	        19.1 ns/op	5019.03 MB/s
BenchmarkNomem/8-4  	50000000	        25.4 ns/op	5033.00 MB/s
BenchmarkNomem/10-4 	50000000	        31.8 ns/op	5024.74 MB/s
BenchmarkNomem/12-4 	50000000	        38.7 ns/op	4956.15 MB/s
BenchmarkNomem/14-4 	30000000	        44.5 ns/op	5029.97 MB/s
PASS
```

## No Memory Constant Block Register

We do see poor performance with a constant block register and varying key
register.

```
$ go version
go version go1.8.1 darwin/amd64
$ git log --oneline -n 1
7a9a7f7 vary key register
$ go test -v -bench .
=== RUN   TestSingle
--- PASS: TestSingle (0.00s)
=== RUN   TestMulti
--- PASS: TestMulti (0.00s)
BenchmarkSingle-4   	200000000	         6.98 ns/op	2293.71 MB/s
BenchmarkMulti/2-4  	200000000	         8.88 ns/op	3604.13 MB/s
BenchmarkMulti/4-4  	100000000	        13.0 ns/op	4937.03 MB/s
BenchmarkMulti/6-4  	100000000	        19.2 ns/op	5002.91 MB/s
BenchmarkMulti/8-4  	50000000	        25.7 ns/op	4983.48 MB/s
BenchmarkMulti/10-4 	50000000	        31.9 ns/op	5019.30 MB/s
BenchmarkMulti/12-4 	50000000	        38.6 ns/op	4977.44 MB/s
BenchmarkMulti/14-4 	30000000	        44.0 ns/op	5089.05 MB/s
BenchmarkNomem/2-4  	50000000	        25.2 ns/op	1270.83 MB/s
BenchmarkNomem/4-4  	30000000	        50.7 ns/op	1262.65 MB/s
BenchmarkNomem/6-4  	20000000	        75.1 ns/op	1278.25 MB/s
BenchmarkNomem/8-4  	20000000	       100 ns/op	1268.01 MB/s
BenchmarkNomem/10-4 	10000000	       126 ns/op	1262.54 MB/s
BenchmarkNomem/12-4 	10000000	       152 ns/op	1260.20 MB/s
BenchmarkNomem/14-4 	10000000	       176 ns/op	1271.54 MB/s
PASS
```

## Longer Tests

Running longer tests (version `9eebcbd`) shows no clear improvement for the no
memory versions against the multi versions. This suggests that the additional
memory and arithmetic operations are essentially free relative to AES-NI.
That's unsurprising but good to confirm.

Our experiments seem to suggest that the single version may be significantly
worse because of its register use. It's not clear how to work around this.

```
$ go test -v -bench . -benchtime 20s
=== RUN   TestSingle
--- PASS: TestSingle (0.00s)
=== RUN   TestMulti
--- PASS: TestMulti (0.00s)
BenchmarkSingle-4   	5000000000	         6.95 ns/op	2303.00 MB/s
BenchmarkMulti/2-4  	5000000000	         7.46 ns/op	4289.18 MB/s
BenchmarkMulti/4-4  	2000000000	        12.8 ns/op	5008.39 MB/s
BenchmarkMulti/6-4  	2000000000	        19.2 ns/op	4989.52 MB/s
BenchmarkMulti/8-4  	1000000000	        25.6 ns/op	4991.35 MB/s
BenchmarkMulti/10-4 	1000000000	        31.9 ns/op	5020.50 MB/s
BenchmarkMulti/12-4 	1000000000	        38.3 ns/op	5015.51 MB/s
BenchmarkMulti/14-4 	1000000000	        44.8 ns/op	4999.47 MB/s
BenchmarkNomem/2-4  	5000000000	         6.39 ns/op	5004.53 MB/s
BenchmarkNomem/4-4  	2000000000	        12.8 ns/op	4986.14 MB/s
BenchmarkNomem/6-4  	2000000000	        19.2 ns/op	4989.02 MB/s
BenchmarkNomem/8-4  	1000000000	        25.6 ns/op	5003.98 MB/s
BenchmarkNomem/10-4 	1000000000	        32.0 ns/op	5007.75 MB/s
BenchmarkNomem/12-4 	1000000000	        38.5 ns/op	4984.65 MB/s
BenchmarkNomem/14-4 	1000000000	        44.7 ns/op	5011.93 MB/s
PASS
```

# stored

#### Benchmarking 

##### safemap
```
BenchmarkInsert-4                2000000               895 ns/op           2.23 MB/s          21 B/op          1 allocs/op
BenchmarkAtomicUpdate-4          1000000              1195 ns/op           1.67 MB/s         122 B/op          3 allocs/op
BenchmarkAtomicWaitUpdate-4      1000000              1981 ns/op           1.01 MB/s         218 B/op          4 allocs/op
BenchmarkUpdate-4                1000000              1191 ns/op           1.68 MB/s          58 B/op          2 allocs/op
BenchmarkAtomicComplex-4         1000000              1301 ns/op           1.54 MB/s         124 B/op          2 allocs/op
BenchmarkAtomicWaitComplex-4     1000000              1802 ns/op           1.11 MB/s         204 B/op          3 allocs/op
BenchmarkAtomicFind-4            2000000              1137 ns/op           1.76 MB/s          96 B/op          2 allocs/op
BenchmarkAtomicWaitFind-4        1000000              1745 ns/op           1.15 MB/s         176 B/op          3 allocs/op
BenchmarkFind-4                  1000000              1728 ns/op           1.16 MB/s         128 B/op          2 allocs/op
BenchmarkEachFullCicle-4             100          15936156 ns/op           0.00 MB/s          64 B/op          1 allocs/op
BenchmarkEachShort-4             2000000              1063 ns/op           1.88 MB/s          64 B/op          1 allocs/op
BenchmarkDelete-4                2000000               607 ns/op           3.29 MB/s           0 B/op          0 allocs/op
ok      github.com/jenchik/stored/safemap 38.704s
```

##### hugemap
```
BenchmarkInsert-4                2000000               984 ns/op           2.03 MB/s          21 B/op          1 allocs/op
BenchmarkAtomicUpdate-4          1000000              2029 ns/op           0.99 MB/s         138 B/op          3 allocs/op
BenchmarkAtomicWaitUpdate-4       500000              3345 ns/op           0.60 MB/s         245 B/op          4 allocs/op
BenchmarkUpdate-4                1000000              1182 ns/op           1.69 MB/s          58 B/op          2 allocs/op
BenchmarkAtomicComplex-4         1000000              1563 ns/op           1.28 MB/s         140 B/op          2 allocs/op
BenchmarkAtomicWaitComplex-4      500000              2532 ns/op           0.79 MB/s         232 B/op          3 allocs/op
BenchmarkAtomicFind-4            1000000              1454 ns/op           1.38 MB/s         112 B/op          2 allocs/op
BenchmarkAtomicWaitFind-4        1000000              1946 ns/op           1.03 MB/s         192 B/op          3 allocs/op
BenchmarkFind-4                  5000000               339 ns/op           5.89 MB/s           0 B/op          0 allocs/op
BenchmarkEachFullCicle-4             100          15072323 ns/op           0.00 MB/s          80 B/op          1 allocs/op
BenchmarkEachShort-4             1000000              1025 ns/op           1.95 MB/s          80 B/op          1 allocs/op
BenchmarkDelete-4                3000000               618 ns/op           3.23 MB/s           0 B/op          0 allocs/op
ok      github.com/jenchik/stored/hugemap 35.077s
```

##### multimap x4
- hugemap
```
BenchmarkInsert-4                        1000000              2407 ns/op           0.83 MB/s         151 B/op          1 allocs/op
BenchmarkAtomicUpdate-4                  1000000              2992 ns/op           0.67 MB/s         158 B/op          3 allocs/op
BenchmarkAtomicWaitUpdate-4               200000              5409 ns/op           0.37 MB/s         325 B/op          3 allocs/op
BenchmarkUpdate-4                        1000000              2718 ns/op           0.74 MB/s         140 B/op          2 allocs/op
BenchmarkAtomicComplex-4                 1000000              2235 ns/op           0.89 MB/s         175 B/op          2 allocs/op
BenchmarkAtomicWaitComplex-4             1000000              1299 ns/op           1.54 MB/s         140 B/op          2 allocs/op
BenchmarkAtomicFind-4                    2000000               660 ns/op           3.03 MB/s          96 B/op          2 allocs/op
BenchmarkAtomicWaitFind-4                2000000               602 ns/op           3.32 MB/s          80 B/op          2 allocs/op
BenchmarkFind-4                          5000000               396 ns/op           5.04 MB/s           0 B/op          0 allocs/op
BenchmarkDelete-4                        1000000              2433 ns/op           0.82 MB/s          18 B/op          0 allocs/op
```
- safemap
```
BenchmarkUnsafeInsert-4                  1000000              1254 ns/op           1.59 MB/s          26 B/op          1 allocs/op
BenchmarkUnsafeAtomicUpdate-4            1000000              1137 ns/op           1.76 MB/s         122 B/op          3 allocs/op
BenchmarkUnsafeAtomicWaitUpdate-4        1000000              1852 ns/op           1.08 MB/s         218 B/op          4 allocs/op
BenchmarkUnsafeUpdate-4                  1000000              1211 ns/op           1.65 MB/s          58 B/op          2 allocs/op
BenchmarkUnsafeAtomicComplex-4           1000000              1212 ns/op           1.65 MB/s         124 B/op          2 allocs/op
BenchmarkUnsafeAtomicWaitComplex-4       1000000              1651 ns/op           1.21 MB/s         204 B/op          3 allocs/op
BenchmarkUnsafeAtomicFind-4              2000000               674 ns/op           2.97 MB/s          96 B/op          2 allocs/op
BenchmarkUnsafeAtomicWaitFind-4          2000000               590 ns/op           3.39 MB/s          80 B/op          2 allocs/op
BenchmarkUnsafeFind-4                    3000000               417 ns/op           4.79 MB/s           0 B/op          0 allocs/op
BenchmarkUnsafeDelete-4                  1000000              2328 ns/op           0.86 MB/s          15 B/op          0 allocs/op
ok      github.com/jenchik/stored/multimap 117.809s
```

### TODO

Testing:

- api.Mapper for all ***map
- Cache, L2Cache, LiveCache, WBCache

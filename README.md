# stored

#### Benchmarking 

##### safemap
```
BenchmarkInsert-4                        2000000              1148 ns/op           1.74 MB/s          21 B/op          1 allocs/op
BenchmarkAtomicUpdate-4                  1000000              1587 ns/op           1.26 MB/s         122 B/op          3 allocs/op
BenchmarkAtomicWaitUpdate-4              1000000              2443 ns/op           0.82 MB/s         202 B/op          4 allocs/op
BenchmarkUpdate-4                        1000000              1441 ns/op           1.39 MB/s          58 B/op          2 allocs/op
BenchmarkAtomicComplex-4                 1000000              1447 ns/op           1.38 MB/s         108 B/op          2 allocs/op
BenchmarkAtomicWaitComplex-4             1000000              2089 ns/op           0.96 MB/s         188 B/op          3 allocs/op
BenchmarkAtomicFind-4                    1000000              1283 ns/op           1.56 MB/s          96 B/op          2 allocs/op
BenchmarkAtomicWaitFind-4                1000000              2070 ns/op           0.97 MB/s         176 B/op          3 allocs/op
BenchmarkFind-4                          1000000              1883 ns/op           1.06 MB/s         128 B/op          2 allocs/op
BenchmarkEachFullCicle-4                     100          26074364 ns/op           0.00 MB/s          64 B/op          1 allocs/op
BenchmarkEachShort-4                     2000000              1234 ns/op           1.62 MB/s          63 B/op          0 allocs/op
BenchmarkDelete-4                        3000000               499 ns/op           4.00 MB/s           0 B/op          0 allocs/op
BenchmarkThreadsInsert-4                 2000000               814 ns/op           2.45 MB/s          16 B/op          1 allocs/op
BenchmarkThreadsAtomicUpdate-4           1000000              1088 ns/op           1.84 MB/s         112 B/op          3 allocs/op
BenchmarkThreadsAtomicWaitUpdate-4        500000              2173 ns/op           0.92 MB/s         193 B/op          4 allocs/op
BenchmarkThreadsUpdate-4                 1000000              1116 ns/op           1.79 MB/s          48 B/op          2 allocs/op
BenchmarkThreadsAtomicComplex-4          2000000              1101 ns/op           1.82 MB/s          96 B/op          2 allocs/op
BenchmarkThreadsAtomicWaitComplex-4      1000000              2325 ns/op           0.86 MB/s         176 B/op          3 allocs/op
BenchmarkThreadsAtomicFind-4             2000000              1120 ns/op           1.79 MB/s          96 B/op          2 allocs/op
BenchmarkThreadsAtomicWaitFind-4         1000000              2291 ns/op           0.87 MB/s         176 B/op          3 allocs/op
BenchmarkThreadsFind-4                   1000000              2317 ns/op           0.86 MB/s         128 B/op          2 allocs/op
BenchmarkThreadsEachShort-4              1000000              1301 ns/op           1.54 MB/s          64 B/op          1 allocs/op
BenchmarkThreadsDelete-4                 2000000               842 ns/op           2.37 MB/s           0 B/op          0 allocs/op
ok      github.com/jenchik/stored/safemap 80.854s
```

##### hugemap
```
BenchmarkInsert-4                        1000000              1064 ns/op           1.88 MB/s          26 B/op          1 allocs/op
BenchmarkAtomicUpdate-4                  1000000              1946 ns/op           1.03 MB/s          58 B/op          2 allocs/op
BenchmarkAtomicWaitUpdate-4              1000000              3069 ns/op           0.65 MB/s         138 B/op          3 allocs/op
BenchmarkUpdate-4                        1000000              1572 ns/op           1.27 MB/s          58 B/op          2 allocs/op
BenchmarkAtomicComplex-4                 1000000              1450 ns/op           1.38 MB/s          44 B/op          1 allocs/op
BenchmarkAtomicWaitComplex-4              500000              2433 ns/op           0.82 MB/s         136 B/op          2 allocs/op
BenchmarkAtomicFind-4                    1000000              1433 ns/op           1.40 MB/s          32 B/op          1 allocs/op
BenchmarkAtomicWaitFind-4                1000000              2164 ns/op           0.92 MB/s         112 B/op          2 allocs/op
BenchmarkFind-4                          3000000               421 ns/op           4.75 MB/s           0 B/op          0 allocs/op
BenchmarkEachFullCicle-4                     100          24757402 ns/op           0.00 MB/s           0 B/op          0 allocs/op
BenchmarkEachShort-4                     2000000              1147 ns/op           1.74 MB/s           0 B/op          0 allocs/op
BenchmarkDelete-4                        2000000               726 ns/op           2.75 MB/s           0 B/op          0 allocs/op
BenchmarkThreadsInsert-4                 2000000               745 ns/op           2.68 MB/s          16 B/op          1 allocs/op
BenchmarkThreadsAtomicUpdate-4           1000000              1473 ns/op           1.36 MB/s          48 B/op          2 allocs/op
BenchmarkThreadsAtomicWaitUpdate-4        500000              3230 ns/op           0.62 MB/s         129 B/op          3 allocs/op
BenchmarkThreadsUpdate-4                 1000000              1004 ns/op           1.99 MB/s          48 B/op          2 allocs/op
BenchmarkThreadsAtomicComplex-4          2000000               921 ns/op           2.17 MB/s          32 B/op          1 allocs/op
BenchmarkThreadsAtomicWaitComplex-4      1000000              2483 ns/op           0.81 MB/s         112 B/op          2 allocs/op
BenchmarkThreadsAtomicFind-4             2000000               912 ns/op           2.19 MB/s          32 B/op          1 allocs/op
BenchmarkThreadsAtomicWaitFind-4         1000000              2707 ns/op           0.74 MB/s         112 B/op          2 allocs/op
BenchmarkThreadsFind-4                  10000000               166 ns/op          11.99 MB/s           0 B/op          0 allocs/op
BenchmarkThreadsEachShort-4              1000000              1050 ns/op           1.90 MB/s           0 B/op          0 allocs/op
BenchmarkThreadsDelete-4                 2000000               831 ns/op           2.41 MB/s           0 B/op          0 allocs/op
ok      github.com/jenchik/stored/hugemap 72.779s
```

##### simplemap
```
BenchmarkInsert-4                        3000000               399 ns/op           5.00 MB/s          19 B/op          1 allocs/op
BenchmarkAtomicUpdate-4                  1000000              1346 ns/op           1.49 MB/s          58 B/op          2 allocs/op
BenchmarkAtomicWaitUpdate-4              1000000              1559 ns/op           1.28 MB/s         106 B/op          3 allocs/op
BenchmarkUpdate-4                        2000000               946 ns/op           2.11 MB/s          53 B/op          2 allocs/op
BenchmarkAtomicComplex-4                 2000000               682 ns/op           2.93 MB/s          38 B/op          1 allocs/op
BenchmarkAtomicWaitComplex-4             2000000               889 ns/op           2.25 MB/s          86 B/op          2 allocs/op
BenchmarkAtomicFind-4                    2000000               660 ns/op           3.03 MB/s          32 B/op          1 allocs/op
BenchmarkAtomicWaitFind-4                2000000               876 ns/op           2.28 MB/s          80 B/op          2 allocs/op
BenchmarkFind-4                          5000000               269 ns/op           7.42 MB/s           0 B/op          0 allocs/op
BenchmarkEachFullCicle-4                      50          23578052 ns/op           0.00 MB/s          32 B/op          1 allocs/op
BenchmarkEachShort-4                     2000000               716 ns/op           2.79 MB/s          32 B/op          1 allocs/op
BenchmarkDelete-4                       10000000               101 ns/op          19.71 MB/s           0 B/op          0 allocs/op
BenchmarkThreadsInsert-4                 2000000               668 ns/op           2.99 MB/s          16 B/op          1 allocs/op
BenchmarkThreadsAtomicUpdate-4           1000000              2003 ns/op           1.00 MB/s          48 B/op          2 allocs/op
BenchmarkThreadsAtomicWaitUpdate-4       1000000              1160 ns/op           1.72 MB/s          97 B/op          3 allocs/op
BenchmarkThreadsUpdate-4                 2000000               812 ns/op           2.46 MB/s          48 B/op          2 allocs/op
BenchmarkThreadsAtomicComplex-4          1000000              1294 ns/op           1.54 MB/s          32 B/op          1 allocs/op
BenchmarkThreadsAtomicWaitComplex-4      5000000               638 ns/op           3.13 MB/s          80 B/op          2 allocs/op
BenchmarkThreadsAtomicFind-4             1000000              1295 ns/op           1.54 MB/s          32 B/op          1 allocs/op
BenchmarkThreadsAtomicWaitFind-4         5000000               648 ns/op           3.08 MB/s          80 B/op          2 allocs/op
BenchmarkThreadsFind-4                  10000000               141 ns/op          14.11 MB/s           0 B/op          0 allocs/op
BenchmarkThreadsEachShort-4              2000000               861 ns/op           2.32 MB/s          32 B/op          1 allocs/op
BenchmarkThreadsDelete-4                 3000000               566 ns/op           3.53 MB/s           0 B/op          0 allocs/op
ok      github.com/jenchik/stored/simplemap       88.559s
```

##### multimap x4:
- hugemap
```
BenchmarkInsert-4                        1000000              2470 ns/op           0.81 MB/s         150 B/op          1 allocs/op
BenchmarkAtomicUpdate-4                  1000000              3081 ns/op           0.65 MB/s         134 B/op          3 allocs/op
BenchmarkAtomicWaitUpdate-4               200000              5757 ns/op           0.35 MB/s         277 B/op          3 allocs/op
BenchmarkUpdate-4                        1000000              2820 ns/op           0.71 MB/s         138 B/op          2 allocs/op
BenchmarkAtomicComplex-4                  300000              3750 ns/op           0.53 MB/s         172 B/op          3 allocs/op
BenchmarkAtomicWaitComplex-4              300000              3653 ns/op           0.55 MB/s         195 B/op          2 allocs/op
BenchmarkAtomicFind-4                    2000000               719 ns/op           2.78 MB/s          64 B/op          2 allocs/op
BenchmarkAtomicWaitFind-4                2000000               649 ns/op           3.08 MB/s          48 B/op          2 allocs/op
BenchmarkFind-4                          3000000               493 ns/op           4.05 MB/s           0 B/op          0 allocs/op
BenchmarkDelete-4                        1000000              2045 ns/op           0.98 MB/s          19 B/op          0 allocs/op
BenchmarkThreadsInsert-4                 1000000              2190 ns/op           0.91 MB/s         471 B/op          2 allocs/op
BenchmarkThreadsDelete-4                  500000              2350 ns/op           0.85 MB/s         444 B/op          1 allocs/op
```
- safemap
```
BenchmarkUnsafeInsert-4                  1000000              2623 ns/op           0.76 MB/s          26 B/op          1 allocs/op
BenchmarkUnsafeAtomicUpdate-4            1000000              1757 ns/op           1.14 MB/s         122 B/op          3 allocs/op
BenchmarkUnsafeAtomicWaitUpdate-4        1000000              2023 ns/op           0.99 MB/s         202 B/op          4 allocs/op
BenchmarkUnsafeUpdate-4                  1000000              1452 ns/op           1.38 MB/s          58 B/op          2 allocs/op
BenchmarkUnsafeAtomicComplex-4           1000000              1425 ns/op           1.40 MB/s         108 B/op          2 allocs/op
BenchmarkUnsafeAtomicWaitComplex-4       1000000              1998 ns/op           1.00 MB/s         188 B/op          3 allocs/op
BenchmarkUnsafeAtomicFind-4              1000000              1169 ns/op           1.71 MB/s          96 B/op          2 allocs/op
BenchmarkUnsafeAtomicWaitFind-4          1000000              2093 ns/op           0.96 MB/s         176 B/op          3 allocs/op
BenchmarkUnsafeFind-4                    1000000              2094 ns/op           0.95 MB/s         128 B/op          2 allocs/op
BenchmarkUnsafeDelete-4                  2000000               972 ns/op           2.06 MB/s           0 B/op          0 allocs/op
ok      github.com/jenchik/stored/multimap        196.197s
```

### TODO

Testing:

- api.Mapper for all ***map
- Cache, L2Cache, LiveCache, WBCache

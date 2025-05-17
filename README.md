# object-pool-go

### 📊 Benchmark Results

| Scenario                             | ns/op | B/op | allocs/op | Description                   |
|--------------------------------------|-------|------|-----------|-------------------------------|
| Benchmark_NewWithoutPool-10          | 130   | 1024 | 1         | 每次都分配新对象                      |
| Benchmark_PoolReuse-10               | 4.192 | 0    | 0         | 非并发池复用对象                      |
| Benchmark_ConcurrencyPoolReuse-10    | 249.9 | 0    | 0         | 并发池，线程安全复用                    |
| Benchmark_NewWithoutPool_Parallel-10 | 248.2 | 1024 | 1         | 并发创建对象（不使用池）                  |
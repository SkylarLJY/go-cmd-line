# Improve CLI tool performance 

`sum()` & `add()` have the same signature: genralize with `statFunc()` to represent both

Perform checks for user experience: give feedbacks to user if they are using tool wrong

Float number comparison is inheretently inprecise: introduce a small tolerence to work around the issue -> customized comparison function 

Comparing complex data structure: helper function or external packages 

Go Benchmark: 
- `testing.B`
- iterate with `b.N` as upper limit: `b.N` is adjusted for the benchmark func to run ~1s
- `go test -bench <bench-regex> -run <test-eacape-regex>`

Profiling:
- to understand how the program spends CPU time
- by adding code to the program: needs maintainace for the additional coe
- by running Go profiler: 
    - `-cpuprofile <output.pprof>` with benchmarks cmd to generate a pprof file
    - `go tool pprof <output.pprof>` to profile
- Memory allocation takes long: read csv files by line instead of read all at once

Tracing
- `trace` with `go test` cmd

`select` block for channels 

## Exercises
Implemented `Min` and `Max` operations

`Min` implemented with CPU numbers of goroutines & `Max` implemented with simple for loop & comparison: due to excess memory allocation from concurrency, a simple implementation like `Max` is more favourable 
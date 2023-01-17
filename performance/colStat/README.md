# Improve CLI tool performance 

`sum()` & `add()` have the same signature: genralize with `statFunc()` to represent both

Perform checks for user experience: give feedbacks to user if they are using tool wrong

Float number comparison is inheretently inprecise: introduce a small tolerence to work around the issue -> customized comparison function 

Comparing complex data structure: helper function or external packages 

Go Benchmark: 
- `testing.B`
- iterate with `b.N` as upper limit: `b.N` is adjusted for the benchmark func to run ~1s
- `go test -bench <bench-regex> -run <test-eacape-regex>`
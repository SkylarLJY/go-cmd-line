# Walk file system
Table driven testing to cover different variations of use cases: have a list of test cases then iterate with a loop

Working with fs: extra caution when deleting files. Never run as a preivileged user. 

Test helpers: when testing file deletion, we change the structure of the fs. Automate test dir creation & cleanup with a test helper function marked with `testing.Helper()` (call in the function).

Good practice for cmd line tools to provide feedbacks to user with stdout: `log.Logger`

`go install` to install the cmd line tool to go bin
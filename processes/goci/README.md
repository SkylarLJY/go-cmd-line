# Executing External Programs
`exec.Command()` 

Execute `go build` w/o creating an executable file: go does not create executable files when building multiple modules -> can add a module from go lib to suppress 

Sometimes a failing program returns a success code and indicates the errors by printing to stdout or stderr: need to capture a program's output to determine the execution correctness 
# Executing External Programs
`exec.Command()` 

Execute `go build` w/o creating an executable file: go does not create executable files when building multiple modules -> can add a module from go lib to suppress 

Sometimes a failing program returns a success code and indicates the errors by printing to stdout or stderr: need to capture a program's output to determine the execution correctness 

Running external commands: some external cmds can potentially hang due to network issues etc -> add timeout

External env changes and to ensure the execution for testing of external commands we can
1. instantiate a local service (eg. a local git server)
2. mock the service: doesn't need to execute the command on the testing machine, can simulate abnormal conditions
    - `TestHelperPrecess()` to replace the regular testing binary
    - `GO_WANT_HELPER_PROCESS` set to 1 to only execute this function as part of the mock test

Handling signals: properly handle signals to clean up resources, save data & exit cleanly on interrupt 

**Bug in example code?** `TestKillRun` testcase 3 with `SIGQUIT` should return timeout error instead of nil as this signal should not be picked up by the program & it will run until timeout. 

Handling json: need to export (capitalize) the field name to unmarshall
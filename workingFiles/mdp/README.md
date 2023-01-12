# Markdown preview as HTML
Assume lib functions are already tested: if you don't trust the lib don't use the lib

Golden files testing: expected outputs saved to files loaded during the test. Good for handling complex outputs. 

The naive implementation adds a `.html` to the end of the md file. This will clutter the dir as a new file is created everytime. It's not thread safe: if two routines try to work on the same file there will be a clash. -> use temp file

cleaning up temp files: 
- `defer os.Remove(<file>)`
- `os.Exit()` exits immediately without excuting deferred statements 
- race condition: what if a file is deleted before the browser has time to open for preview: add delay after executing the preview command

Testing stdout outputs: pass `io.Writer` interface as a parameter, set to stdout for regular functions and buffer for testing 

## Templates

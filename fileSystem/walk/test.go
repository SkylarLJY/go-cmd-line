package main

import "os"

func TestOpen() {
	path := "/var/folders/jz/_w9fhfvd0z5cwrkkp0363lr80000gn/T/walktest3328884551/file1.vim"
	err := os.WriteFile(path, []byte("test"), 0644)
	if err != nil {
		panic(err)
	}
}

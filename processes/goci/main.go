package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type executer interface {
	execute() (string, error)
}

func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("project directory is required:%w", ErrValidation)
	}

	pipeline := make([]executer, 3)
	pipeline[0] = newStep("go build", "go", "Go build: SUCCESS", proj, []string{"build", ".", "errors"})
	pipeline[1] = newStep("go test", "go", "Go test: SUCCESS", proj, []string{"test", "-v"})
	pipeline[2] = newExceptionStep("go fmt", "gofmt", "Go fmt: SUCCESS", proj, []string{"-l", "."})
	for _, s := range pipeline {
		msg, err := s.execute()
		if err != nil {
			return err
		}
		fmt.Fprintln(out, msg)
	}

	return nil
}

func main() {
	proj := flag.String("p", "", "project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

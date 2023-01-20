package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("project directory is required:%w", ErrValidation)
	}

	args := []string{"build", ".", "errors"} // add error package so that it doesn't generate an executable file
	cmd := exec.Command("go", args...)
	cmd.Dir = proj
	if err := cmd.Run(); err != nil {
		return &stepErr{"go build", "go build failed", err}
	}
	_, err := fmt.Fprintln(out, "Go build: SUCCESS")
	return err
}

func main() {
	proj := flag.String("p", "", "project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

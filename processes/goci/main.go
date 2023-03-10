package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

type executer interface {
	execute() (string, error)
}

// var pipeline []executer

func run(proj, branch, pipelineFile string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("project directory is required:%w", ErrValidation)
	}

	pipeline, err := initPipeline(pipelineFile, proj, branch)
	if err != nil {
		return err
	}

	// pipeline := make([]executer, 6)
	// pipeline[0] = newStep("go build", "go", "Go build: SUCCESS", proj, []string{"build", ".", "errors"})
	// pipeline[1] = newStep("go test", "go", "Go test: SUCCESS", proj, []string{"test", "-v"})
	// pipeline[2] = newExceptionStep("go fmt", "gofmt", "Go fmt: SUCCESS", proj, []string{"-l", "."})
	// pipeline[3] = newTimeoutStep("git push", "git", "Git push: SUCCESS", proj, []string{"push", "origin", "main"}, 10*time.Second)
	// pipeline[4] = newStep("golangci-lint", "golangci-lint", "Golang Lint: SUCCESS", proj, []string{"run"})
	// pipeline[5] = newPrintStep("gocyclo", "gocyclo", "Gocyclo: SUCCESS", proj, []string{"."})

	sig := make(chan os.Signal, 1)
	errCh := make(chan error)
	done := make(chan struct{})

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	if err := gitCredCheck(proj); err != nil {
		return err
	}

	go func() {
		for _, s := range pipeline {
			msg, err := s.execute()
			if err != nil {
				errCh <- err
			}
			_, err = fmt.Fprintln(out, msg)
			if err != nil {
				errCh <- err
			}

		}
		close(done)
	}()

	for {
		select {
		case rec := <-sig:
			signal.Stop(sig)
			return fmt.Errorf("%s: Exiting: %w", rec, ErrSignal)
		case err := <-errCh:
			return err
		case <-done:
			return nil
		}
	}
}

func main() {
	proj := flag.String("p", "", "project directory")
	branch := flag.String("b", "main", "target brance")
	pipelineFile := flag.String("pipeline", "pipeline.json", "pipeline config")
	flag.Parse()

	// if *branch != "" {
	// 	if err := gitSwitchBranch(*branch, *proj); err != nil {
	// 		fmt.Fprintln(os.Stderr, err)
	// 		os.Exit(1)
	// 	}
	// }

	// pl, err := initPipeline(*pipelineFile, *proj)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	os.Exit(1)
	// }

	// pipeline = pl

	if err := run(*proj, *branch, *pipelineFile, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

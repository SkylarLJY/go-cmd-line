package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

func run(files []string, op string, col int, out io.Writer) error {
	var opFunc statFunc
	// check input validity
	if len(files) == 0 {
		return ErrNoFiles
	}
	if col < 1 {
		return ErrInvalidCol
	}
	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	case "min":
		opFunc = min
	case "max":
		opFunc = max
	default:
		return ErrInvalidOp
	}

	consolidate := make([]float64, 0)
	resCh := make(chan []float64)
	errCh := make(chan error)
	doneCh := make(chan struct{}) // using empty struct b/c we don't need to send any data, just act as a signal
	filesCh := make(chan string)

	go func() {
		defer close(filesCh)
		for _, fname := range files {
			filesCh <- fname
		}
	}()

	wg := sync.WaitGroup{}

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fname := range filesCh {
				f, err := os.Open(fname)
				if err != nil {
					errCh <- fmt.Errorf("cannot open file: %w", err)
					return
				}
				data, err := csv2float(f, col)
				if err != nil {
					errCh <- err
				}
				if err := f.Close(); err != nil {
					errCh <- err
				}
				resCh <- data
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()
	for {
		select {
		case err := <-errCh:
			return err
		case data := <-resCh:
			consolidate = append(consolidate, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(out, opFunc(consolidate))
			return err

		}
	}
}

func main() {
	op := flag.String("op", "sum", "operation to execute")
	col := flag.Int("col", 1, "the column to operate on")
	flag.Parse()

	if err := run(flag.Args(), *op, *col, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

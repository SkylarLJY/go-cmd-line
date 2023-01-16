package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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
	default:
		return ErrInvalidOp
	}

	consolidate := make([]float64, 0)
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()
		data, err := csv2float(f, col)
		if err != nil {
			return err
		}
		consolidate = append(consolidate, data...)
	}

	_, err := fmt.Fprintln(out, opFunc(consolidate))

	return err
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

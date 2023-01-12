package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func count(r io.Reader, countLines bool, countBytes bool) int {
	scanner := bufio.NewScanner(r)
	if countBytes {
		scanner.Split(bufio.ScanBytes)
	} else if !countLines {
		scanner.Split(bufio.ScanWords)
	}
	wc := 0
	for scanner.Scan() {
		wc++
	}
	return wc
}

func main() {
	lines := flag.Bool("l", false, "count lines")
	numBytes := flag.Bool("b", false, "count bytes")
	inputFile := flag.String("f", "", "read input from file")
	flag.Parse()
	input := os.Stdin
	if *inputFile != "" {
		var err error
		input, err = os.Open(*inputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer input.Close()
	}
	fmt.Println(count(input, *lines, *numBytes))

}

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
	inputFiles := flag.Bool("f", false, "read input from files")
	flag.Parse()

	input := os.Stdin
	var res int
	if *inputFiles {
		files := flag.Args()
		var err error
		for _, file := range files {
			input, err = os.Open(file)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer input.Close()

			res += count(input, *lines, *numBytes)
		}
	} else {
		res = count(input, *lines, *numBytes)
	}
	fmt.Println(res)

}

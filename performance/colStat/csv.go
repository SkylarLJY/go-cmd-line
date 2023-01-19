package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"sync"
)

type statFunc func(data []float64) float64

func sum(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

func min(data []float64) float64 {
	numCh := make(chan float64)
	resCh := make(chan float64)
	doneCh := make(chan struct{})
	res := data[0]
	go func() {
		defer close(numCh)
		var n float64
		for _, n = range data {
			numCh <- n
		}
	}()

	wg := sync.WaitGroup{}

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := res
			for num := range numCh {
				// result = math.Min(result, num)
				if result < num {
					result = num
				}
			}
			resCh <- result
		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case newRes := <-resCh:
			// res = math.Min(newRes, res)
			if res > newRes {
				res = newRes
			}
		case <-doneCh:
			return res
		}
	}

	// minVal := data[0]
	// for _, v := range data {
	// 	minVal = math.Min(minVal, v)
	// }
	// return minVal
}

func max(data []float64) float64 {
	maxVal := data[0]
	// var v float64
	for _, v := range data {
		// maxVal = math.Max(maxVal, v)
		if maxVal < v {
			maxVal = v
		}
	}
	return maxVal
}

func csv2float(r io.Reader, col int) ([]float64, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true // reduce mem alloc
	col--

	var data []float64
	for i := 0; ; i++ {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Cannot read from csv file %w", err)
		}
		if i == 0 {
			continue
		}
		if len(row) <= col {
			return nil, fmt.Errorf("%w: file has only %d cols\n", ErrInvalidCol, len(row))
		}
		v, err := strconv.ParseFloat(row[col], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}
		data = append(data, v)
	}
	return data, nil
}

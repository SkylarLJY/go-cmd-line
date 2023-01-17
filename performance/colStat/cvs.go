package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
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

func csv2float(r io.Reader, col int) ([]float64, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true // reduce mem alloc
	col--

	// allData, err := cr.ReadAll()
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot read data from file: %w\n", err)
	// }

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

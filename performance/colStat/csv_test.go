package main

import (
	"bytes"
	"errors"
	"io"
	"math"
	"testing"
	"testing/iotest"
)

const FLOAT_THRESHOLD = 1e10

func floatEqual(a, b, threshold float64) bool {
	return math.Abs(a-b) <= threshold
}

func TestOps(t *testing.T) {
	data := [][]float64{
		{1, 2, 3},
		{.1, .2, .3},
		{111, 90, 87, 102, 444},
	}

	testcases := []struct {
		name     string
		op       statFunc
		expected []float64
	}{
		{"Sum", sum, []float64{6, .6, 834}},
		{"Avg", avg, []float64{2, .2, 166.8}},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < len(data); i++ {
				res := tc.op(data[i])
				if !floatEqual(res, tc.expected[i], FLOAT_THRESHOLD) {
					t.Errorf("expected %f but got %f\n", tc.expected[i], res)
				}
			}
		})
	}
}

func TestCSV2Float(t *testing.T) {
	csvData := `IP, Req, Res Time
192.168.0.199,2056,236
192.168.0.88,899,220
192.168.0.199,3054,226
192.168.0.100,4133,218
192.168.0.199,950,238
`
	testcases := []struct {
		name   string
		col    int
		exp    []float64
		expErr error
		r      io.Reader
	}{
		{"Column2", 2, []float64{2506, 899, 3054, 4133, 950}, nil, bytes.NewBuffer([]byte(csvData))},
		{"Column3", 3, []float64{236, 220, 226, 218, 236}, nil, bytes.NewBuffer([]byte(csvData))},
		{"FailRead", 1, nil, iotest.ErrTimeout, iotest.TimeoutReader(bytes.NewReader([]byte{0}))},
		{"FailedNotNum", 1, nil, ErrNotNumber, bytes.NewBuffer([]byte(csvData))},
		{"FailedInvalidCol", 4, nil, ErrInvalidCol, bytes.NewBuffer([]byte(csvData))},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := csv2float(tc.r, tc.col)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected error %q, but got %q", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}
			for i, exp := range tc.exp {
				if !floatEqual(exp, res[i], FLOAT_THRESHOLD) {
					t.Errorf("expected %f but got %f", exp, res[i])
				}
			}
		})
	}
}

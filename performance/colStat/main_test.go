package main

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	testcase := []struct {
		name   string
		col    int
		op     string
		exp    string
		files  []string
		expErr error
	}{
		{"RunAvg1File", 3, "avg", "227.6\n", []string{"./testdata/example.csv"}, nil},
		{"RunAvgMultiFiles", 3, "avg", "233.84\n", []string{"./testdata/example.csv", "./testdata/example2.csv"}, nil},
		{"RunFailedRead", 3, "avg", "", []string{"./testdata/fakefile.csv"}, os.ErrNotExist},
		{"RunFailCol", 0, "avg", "", []string{"./testdata/example.csv"}, ErrInvalidCol},
		{"RunFailNoFiles", 2, "avg", "", []string{}, ErrNoFiles},
		{"RunFailedOp", 2, "invalidOp", "", []string{"./testdata/example.csv"}, ErrInvalidOp},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			err := run(tc.files, tc.op, tc.col, &buffer)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected error %q but got %q", tc.expErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error %q", err)
			}

			if tc.exp != buffer.String() {
				t.Errorf("expected %q but got %q", tc.exp, &buffer)
			}
		})
	}
}

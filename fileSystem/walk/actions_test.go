package main

import (
	"os"
	"testing"
)

func TestFilterOut(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		ext      string
		minSize  int64
		expected bool
	}{
		{"FilterNoExt", "testdata/dir.log", "", 0, false},
		{"FilterExtMatch", "testdata/dir.log", ".log", 0, false},
		{"FilterExtNoMatch", "testdata/dir.log", ".sh", 0, true},
		{"FilterExtSizeMatch", "testdata/dir.log", ".log", 10, false},
		{"FilterExtSizeNoMatch", "testdata/dir.log", ".log", 100, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := os.Stat(tc.file)
			if err != nil {
				t.Fatal(err)
			}
			f := filterOut(tc.file, tc.ext, tc.minSize, info)
			if f != tc.expected {
				t.Errorf("test case %s: expected %t but got %t\n", tc.name, tc.expected, f)
			}
		})
	}
}

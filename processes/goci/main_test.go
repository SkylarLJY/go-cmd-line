package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	testcase := []struct {
		name   string
		runDir string
		out    string
		expErr error
	}{
		{"success", "./testdata/tool",
			"Go build: SUCCESS\nGo test: SUCCESS\nGo fmt: SUCCESS\nGit push: SUCCESS\n",
			nil},
		{"fail", "./testdata/toolErr", "", &stepErr{step: "go build"}},
		{"failfmt", "./testdata/toolFmtErr", "", &stepErr{step: "go fmt"}},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			err := run(tc.runDir, &buffer)
			if tc.expErr != nil {
				if err == nil || !errors.Is(tc.expErr, err) {
					t.Errorf("Expected error %q but got %q", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %q", err)
				return
			}
			if buffer.String() != tc.out {
				t.Errorf("Expected output %q but got %q", tc.out, buffer.String())
			}

		})
	}
}

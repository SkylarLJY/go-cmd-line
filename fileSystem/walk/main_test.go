package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	for k, n := range files {
		for j := 1; j <= n; j++ {
			fname := fmt.Sprintf("file%d%s", j, k)
			fpath := filepath.Join(tempDir, fname)
			if err := os.WriteFile(fpath, []byte("dummy content"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return tempDir, func() { os.RemoveAll(tempDir) }
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		// test filtering
		{"NoFilter", "testdata", config{"", 0, true, false}, "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{"FilterExtMatch", "testdata", config{".log", 0, true, false}, "testdata/dir.log\n"},
		{"FilterExtNoMatch", "testdata", config{".vim", 0, true, false}, ""},
		{"FilterExtSizeMatch", "testdata", config{".log", 10, true, false}, "testdata/dir.log\n"},
		{"FilterExtSizeNoMatch", "testdata", config{".log", 100, true, false}, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			if err := run(tc.root, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}
			res := buffer.String()
			if res != tc.expected {
				t.Errorf("test case %s: expected %q but got %q\n", tc.name, tc.expected, res)
			}
		})
	}
}

func TestRunDelExt(t *testing.T) {
	testCases := []struct {
		name     string
		cfg      config
		extNoDel string
		nDel     int
		nNoDel   int
		expected string
	}{
		{"DelExtNoMatch", config{ext: ".log", del: true}, ".vim", 0, 10, ""},
		// {"DelExtMatch", config{ext: ".log", del: true}, "", 10, 0, ""},
		// {"DelExtMixed", config{ext: ".log", del: true}, ".gz", 5, 5, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:  tc.nDel,
				tc.extNoDel: tc.nNoDel,
			})
			defer cleanup()

			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}
			res := buffer.String()
			if res != tc.expected {
				t.Errorf("expected %q but got %q\n", tc.expected, res)
			}

			filesLeft, err := os.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}
			if len(filesLeft) != tc.nNoDel {
				t.Errorf("expected %d files left but have %d files left\n", tc.nNoDel, len(filesLeft))
			}
		})
	}
}

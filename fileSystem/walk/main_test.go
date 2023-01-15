package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
		{"NoFilter", "testdata", config{"", 0, true, false, nil, ""}, "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{"FilterExtMatch", "testdata", config{".log", 0, true, false, nil, ""}, "testdata/dir.log\n"},
		{"FilterExtNoMatch", "testdata", config{".vim", 0, true, false, nil, ""}, ""},
		{"FilterExtSizeMatch", "testdata", config{".log", 10, true, false, nil, ""}, "testdata/dir.log\n"},
		{"FilterExtSizeNoMatch", "testdata", config{".log", 100, true, false, nil, ""}, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)
			tc.cfg.wLog = &logBuffer
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
		{"DelExtMatch", config{ext: ".log", del: true}, "", 10, 0, ""},
		{"DelExtMixed", config{ext: ".log", del: true}, ".gz", 5, 5, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)
			tc.cfg.wLog = &logBuffer
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

			// check delete log
			expLogLines := tc.nDel + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(lines) != expLogLines {
				t.Errorf("Expected %d lines of log but got %d\n", expLogLines, len(lines))
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

func TestArchiveFile(t *testing.T) {
	testCases := []struct {
		name         string
		cfg          config
		extNoArchive string
		nArchive     int
		nNoArchive   int
	}{
		{"ArchiveExtNoMatch", config{ext: ".log"}, ".vim", 0, 3},
		{"ArchiveExtMatch", config{ext: ".log"}, ".vim", 3, 0},
		{"ArchiveExtMixed", config{ext: ".log"}, ".vim", 3, 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:      tc.nArchive,
				tc.extNoArchive: tc.nNoArchive,
			})
			defer cleanup()

			archiveDir, cleanupArchive := createTempDir(t, nil) // create an empty dir for archived files
			defer cleanupArchive()

			tc.cfg.archive = archiveDir
			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			pattern := filepath.Join(tempDir, fmt.Sprintf("*%s", tc.cfg.ext))
			expFile, err := filepath.Glob(pattern)
			if err != nil {
				t.Fatal(err)
			}

			expOut := strings.Join(expFile, "\n")

			// remove the last new line in output buffer before comparison
			res := strings.TrimSpace(buffer.String())

			if res != expOut {
				t.Errorf("expected %q but got %q\n", expOut, res)
			}

			filesArchived, err := os.ReadDir(archiveDir)
			if err != nil {
				t.Fatal(err)
			}
			if len(filesArchived) != tc.nArchive {
				t.Errorf("expected to have archived %d files but got %d\n", tc.nArchive, len(filesArchived))
			}
		})
	}
}

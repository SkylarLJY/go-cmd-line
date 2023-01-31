package scan_test

import (
	"errors"
	"go-cmd-line/cobra/pScan/scan"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	testCase := []struct {
		name   string
		host   string
		expLen int
		expErr error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 1, scan.ErrExists},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostList{}
			if err := hl.Add("host1"); err != nil {
				t.Fatal(err)
			}
			err := hl.Add(tc.host)
			if tc.expErr != nil {
				if !errors.Is(tc.expErr, err) {
					t.Fatalf("expected %q but got %q", tc.expErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error %q", err)
			}
			if tc.expLen != len(hl.Hosts) {
				t.Errorf("expected length to be %d but got %d", tc.expLen, len(hl.Hosts))
			}

			if hl.Hosts[1] != tc.host {
				t.Errorf("expected host %q to be at index 1 but got %q", tc.host, hl.Hosts[1])
			}
		})
	}
}

func TestRemove(t *testing.T) {
	testCase := []struct {
		name   string
		host   string
		expLen int
		expErr error
	}{
		{"RemoveExisting", "host1", 0, nil},
		{"RemoveNotFound", "host2", 1, scan.ErrNotExits},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostList{}
			if err := hl.Add("host1"); err != nil {
				t.Fatal(err)
			}
			err := hl.Remove(tc.host)
			if tc.expErr != nil {
				if !errors.Is(err, tc.expErr) {
					t.Fatalf("expected error %q but got %q", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error %q", err)
			}
			if tc.expLen != len(hl.Hosts) {
				t.Errorf("expected length to be %d but got %d", tc.expLen, len(hl.Hosts))
			}
		})
	}
}

func TestSaveLoad(t *testing.T) {
	hl1 := scan.HostList{}
	hl2 := scan.HostList{}

	hostName := "host1"
	hl1.Add(hostName)
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("error creating temp file: %s", err)
	}
	defer os.Remove(f.Name())
	if err := hl1.Save(f.Name()); err != nil {
		t.Fatalf("error savign host list to file: %s", err)
	}
	if err := hl2.Load(f.Name()); err != nil {
		t.Fatalf("error loafing list from file: %s", err)
	}
	if hl1.Hosts[0] != hl2.Hosts[0] {
		t.Errorf("host mismatch")
	}
}

func TestLoadNoFile(t *testing.T) {
	tf, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("failed to create temp dir: %q", err)
	}
	if err := os.Remove(tf); err != nil {
		t.Fatalf("failed to remove temp file")
	}
	hl := &scan.HostList{}
	if err := hl.Load(tf); err != nil {
		t.Errorf("unexpected error %s", err)
	}
}

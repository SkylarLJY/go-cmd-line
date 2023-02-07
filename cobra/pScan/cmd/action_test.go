package cmd

import (
	"bytes"
	"fmt"
	"go-cmd-line/cobra/pScan/scan"
	"io"
	"os"
	"strings"
	"testing"
)

func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	tf, err := os.CreateTemp("", "pScan")
	if err != nil {
		t.Fatal(err)
	}
	tf.Close()

	if initList {
		hl := &scan.HostList{}
		for _, h := range hosts {
			hl.Add(h)
		}
		if err := hl.Save(tf.Name()); err != nil {
			t.Fatal(err)
		}
	}

	return tf.Name(), func() { os.Remove(tf.Name()) }
}

func TestHostActions(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}
	testcases := []struct {
		name     string
		args     []string
		expOut   string
		initList bool
		action   func(io.Writer, string, []string) error
	}{
		{"AddAction", hosts, "Added host: host1\nAdded host: host2\nAdded host: host3\n", false, addAction},
		{"ListAction", []string{}, "host1\nhost2\nhost3\n", true, listAction},
		{"DeleteAction", []string{"host1"}, "Deleted host: host1\n", true, deleteAction},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			tf, cleanup := setup(t, hosts, tc.initList)
			defer cleanup()
			if err := tc.action(&out, tf, tc.args); err != nil {
				t.Fatalf("unexpected error %q\n", err)
			}
			if out.String() != tc.expOut {
				t.Errorf("expected %q but got %q", tc.expOut, out.String())
			}
		})
	}
}

func TestIntegration(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}
	tf, cleanup := setup(t, hosts, false)
	defer cleanup()
	delHost := "host2"
	endHosts := []string{"host1", "host3"}
	var out bytes.Buffer
	expOut := ""
	// form add
	for _, v := range hosts {
		expOut += fmt.Sprintf("Added host: %s\n", v)
	}
	// from list
	expOut += strings.Join(hosts, "\n")
	expOut += "\n"
	//from delete
	expOut += fmt.Sprintf("Deleted host: %s\n", delHost)
	// list again
	expOut += strings.Join(endHosts, "\n")
	expOut += "\n"

	if err := addAction(&out, tf, hosts); err != nil {
		t.Fatal(err)
	}

	if err := listAction(&out, tf, nil); err != nil {
		t.Fatal(err)
	}

	if err := deleteAction(&out, tf, []string{delHost}); err != nil {
		t.Fatal(err)
	}

	if err := listAction(&out, tf, nil); err != nil {
		t.Fatal(err)
	}

	if out.String() != expOut {
		t.Errorf("expected output %q but got %q\n", expOut, out.String())
	}

}

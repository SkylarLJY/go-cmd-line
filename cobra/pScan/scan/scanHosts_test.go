package scan_test

import (
	"go-cmd-line/cobra/pScan/scan"
	"net"
	"strconv"
	"testing"
)

func TestStateString(t *testing.T) {
	ps := scan.PortState{}
	if ps.Open.String() != "closed" {
		t.Errorf("expected %q but got %q\n", "closed", ps.Open.String())
	}
	ps.Open = true
	if ps.Open.String() != "open" {
		t.Errorf("expected %q but got %q\n", "open", ps.Open.String())
	}
}

func TestRunHostFound(t *testing.T) {
	testcases := []struct {
		name     string
		expState string
	}{
		{"OpenPort", "open"},
		{"ClosedPort", "closed"},
	}

	host := "127.0.0.1"
	hl := &scan.HostList{}
	hl.Add(host)
	ports := []int{}

	for _, tc := range testcases {
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		ports = append(ports, port)

		if tc.name == "ClosedPort" {
			ln.Close()
		}
	}

	res := scan.Run(hl, ports)
	if len(res) != 1 {
		t.Fatalf("expected 1 result but got %d\n", len(res))
	}
	if res[0].Host != host {
		t.Errorf("expected host %q but got %q\n", host, res[0].Host)
	}
	if res[0].NotFound {
		t.Errorf("%q expected but not found\n", host)
	}

	if len(res[0].PortStates) != 2 {
		t.Fatalf("expected 2 ports but got %d\n", len(res[0].PortStates))
	}

	for i, tc := range testcases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("expected port %d but got port %d\n", ports[i], res[0].PortStates[i].Port)
		}
		if res[0].PortStates[i].Open.String() != tc.expState {
			t.Errorf("expected port %d to be %s\n", ports[i], tc.expState)
		}
	}
}

func TestHostNotFound(t *testing.T) {
	host := "444.555.666.777"
	hl := &scan.HostList{}
	hl.Add(host)

	res := scan.Run(hl, []int{})
	if len(res) != 1 {
		t.Fatalf("expected 1 result but got %d\n", len(res))
	}
	if res[0].Host != host {
		t.Fatalf("expected host %q but got %q\n", host, res[0].Host)
	}
	if !res[0].NotFound {
		t.Errorf("expected host %q NOT to be found\n", host)
	}
	if len(res[0].PortStates) != 0 {
		t.Errorf("expected 0 port state but got %d\n", len(res[0].PortStates))
	}
}

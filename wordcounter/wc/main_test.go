package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("w1 w2 w3 w4\n")
	exp := 4
	res := count(b, false, false)
	if res != exp {
		t.Errorf("expected %d but got %d\n", exp, res)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("w1 w2\nw3 w4\nw5")
	exp := 3
	res := count(b, true, false)
	if res != exp {
		t.Errorf("expected %d but got %d\n", exp, res)
	}
}

func TestCountBytes(t *testing.T) {
	b := bytes.NewBufferString("w1 w2 w3\n")
	exp := 9
	res := count(b, false, true)
	if res != exp {
		t.Errorf("expected %d but got %d\n", exp, res)
	}
}

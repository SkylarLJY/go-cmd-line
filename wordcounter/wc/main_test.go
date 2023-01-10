package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("w1 w2 w3 w4\n")
	exp := 4
	res := count(b)
	if res != exp {
		t.Errorf("expected %d but got %d\n", exp, res)
	}
}

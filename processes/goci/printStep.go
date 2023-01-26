package main

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

type printStep struct {
	step
}

func newPrintStep(name, exe, message, proj string, args []string) printStep {
	s := printStep{}
	s.step = newStep(name, exe, message, proj, args)

	return s
}

func (s printStep) execute() (string, error) {
	var out bytes.Buffer
	cmd := exec.Command(s.exe, s.args...)
	cmd.Dir = s.proj
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	lines := strings.Split(out.String(), "\n")
	for _, l := range lines {
		complexityStr := strings.Split(l, " ")[0]
		complexity, _ := strconv.Atoi(complexityStr)
		if complexity >= 10 {
			return "", ErrHighComplecity
		}

	}

	return s.message, nil
}

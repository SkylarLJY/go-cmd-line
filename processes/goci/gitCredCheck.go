package main

import (
	"bytes"
	"os"
	"os/exec"
)

func gitCredCheck(dir string) error {
	envMapping := map[string]string{
		"user.name":  "GIT_USER_NAME",
		"user.email": "GIT_USER_EMAIL",
	}
	for field, envName := range envMapping {
		var out bytes.Buffer
		cmd := exec.Command("git", "config", field)
		cmd.Dir = dir
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			return err
		}
		if out.Len() == 1 {
			val := os.Getenv(envName)
			if val == "" {
				return ErrGitCredMissing
			}
			cmd := exec.Command("git", "config", field, val)
			cmd.Dir = dir
			if err := cmd.Run(); err != nil {
				return err
			}

		}
	}
	return nil
}

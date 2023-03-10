package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func setupGit(t *testing.T, proj string) func() {
	t.Helper()
	// check if git is available
	gitExec, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	tempDir, err := os.MkdirTemp("", "gocitest")
	if err != nil {
		t.Fatal(err)
	}
	projPath, err := filepath.Abs(proj)
	if err != nil {
		t.Fatal(err)
	}
	remoteURI := fmt.Sprintf("file://%s", tempDir)

	gitCmds := []struct {
		args []string
		dir  string
		env  []string
	}{
		{[]string{"init", "--bare"}, tempDir, nil},
		{[]string{"init"}, projPath, nil},
		{[]string{"remote", "add", "origin", remoteURI}, projPath, nil},
		{[]string{"add", "."}, projPath, nil},
		{[]string{"commit", "-m", "test"}, projPath, []string{
			"GIT_COMMITTER_NAME=test",
			"GIT_COMMITTER_EMAIL=test@example.com",
			"GIT_AUTHOR_NAME=test",
			"GIT_AUTHOR_EMAIL=test@example.com",
		}},
	}

	for _, g := range gitCmds {
		gitCmd := exec.Command(gitExec, g.args...)
		gitCmd.Dir = g.dir
		if g.env != nil {
			gitCmd.Env = append(gitCmd.Env, g.env...)
		}

		if err := gitCmd.Run(); err != nil {
			t.Fatal(err)
		}
	}

	return func() {
		os.RemoveAll(tempDir)
		os.RemoveAll(filepath.Join(projPath, ".git"))
	}
}

func mockCmdContext(ctx context.Context, exe string, args ...string) *exec.Cmd {
	testArgs := []string{"-test.run=TestHelperProcess"}
	testArgs = append(testArgs, exe)
	testArgs = append(testArgs, args...)
	cmd := exec.CommandContext(ctx, os.Args[0], testArgs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func mockCmdTimeout(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cmd := mockCmdContext(ctx, exe, args...)
	cmd.Env = append(cmd.Env, "GO_HELPER_TIMEOUT=1")
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if os.Getenv("GO_HELPER_TIMEOUT") == "1" {
		time.Sleep(15 * time.Second)
	}

	if os.Args[2] == "git" {
		fmt.Fprintln(os.Stdout, "Everything up-to-date")
		os.Exit(0)
	}
	os.Exit(1)
}

func TestRun(t *testing.T) {
	testcase := []struct {
		name     string
		runDir   string
		out      string
		expErr   error
		setupGit bool
		mockCmd  func(ctx context.Context, name string, args ...string) *exec.Cmd
	}{
		{"success", "./testdata/tool",
			"Go build: SUCCESS\nGo test: SUCCESS\nGo fmt: SUCCESS\nGit push: SUCCESS\nGolang Lint: SUCCESS\nGocyclo: SUCCESS\n",
			nil, true, nil},
		{"successMock", "./testdata/tool",
			"Go build: SUCCESS\nGo test: SUCCESS\nGo fmt: SUCCESS\nGit push: SUCCESS\nGolang Lint: SUCCESS\nGocyclo: SUCCESS\n",
			nil, false, mockCmdContext},
		{"fail", "./testdata/toolErr", "", &stepErr{step: "go build"}, false, nil},
		{"failfmt", "./testdata/toolFmtErr", "", &stepErr{step: "go fmt"}, false, nil},
		{"failTimeout", "./testdata/tool", "", context.DeadlineExceeded, false, mockCmdTimeout},
		{"failLint", "./testdata/toolLintErr", "", &stepErr{step: "golangci-lint"}, false, mockCmdContext},
		{"failCyclo", "./testdata/toolErrCyclo", "", ErrHighComplecity, false, mockCmdContext},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupGit {
				_, err := exec.LookPath("git")
				if err != nil {
					t.Skip("Git not installed. Skipping test.")
				}
				cleanup := setupGit(t, tc.runDir)
				defer cleanup()
			}

			if tc.mockCmd != nil {
				command = tc.mockCmd
			}

			var buffer bytes.Buffer
			err := run(tc.runDir, "main", "./pipeline.json", &buffer)
			if tc.expErr != nil {
				if err == nil || !errors.Is(err, tc.expErr) {
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

func TestRunKill(t *testing.T) {
	testCases := []struct {
		name   string
		proj   string
		sig    syscall.Signal
		expErr error
	}{
		{"SIGINT", "./testdata/tool", syscall.SIGINT, ErrSignal},
		{"SIGTERM", "./testdata/tool", syscall.SIGTERM, ErrSignal},
		{"SIGQUIT", "./testdata/tool", syscall.SIGQUIT, &stepErr{step: "git push"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			command = mockCmdTimeout
			errCh := make(chan error, 1)
			ignSigCh := make(chan os.Signal, 1)
			expSigCh := make(chan os.Signal, 1)

			signal.Notify(ignSigCh, syscall.SIGQUIT)
			defer signal.Stop(ignSigCh)

			signal.Notify(expSigCh, tc.sig)
			defer signal.Stop(expSigCh)

			go func() {
				errCh <- run(tc.proj, "main", "pipeline.json", io.Discard)
			}()

			go func() {
				time.Sleep(2 * time.Second)
				syscall.Kill(syscall.Getpid(), tc.sig)
			}()

			err := <-errCh
			if !errors.Is(err, tc.expErr) {
				t.Errorf("expected %q but got %q", tc.expErr, err)
			}

			select {
			case rec := <-expSigCh:
				if rec != tc.sig {
					t.Errorf("expected signal %q but got %q", tc.sig, rec)
				}
			case <-ignSigCh:
			default:
				t.Errorf("no signal received")

			}

		})
	}

}

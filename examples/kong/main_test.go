package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
)

// runCLI mimics main(): parse args, print debug if set, then run.
func runCLI(args []string) (stdout, stderr string, err error) {
	var cli CLI
	k, err := kong.New(&cli,
		kong.Name("clu"),
		kong.Description("Kong example with facts/collector/requires subcommands."),
		kong.UsageOnError(),
	)
	if err != nil {
		return "", "", err
	}
	ctx, err := k.Parse(args)
	if err != nil {
		return "", "", err
	}

	// Capture stdout/stderr
	origOut, origErr := os.Stdout, os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout, os.Stderr = wOut, wErr
	defer func() {
		os.Stdout, os.Stderr = origOut, origErr
	}()

	var bufOut, bufErr bytes.Buffer
	doneOut := make(chan struct{})
	doneErr := make(chan struct{})
	go func() { _, _ = io.Copy(&bufOut, rOut); close(doneOut) }()
	go func() { _, _ = io.Copy(&bufErr, rErr); close(doneErr) }()

	if cli.Debug {
		fmt.Fprintln(os.Stderr, "debug: enabled")
	}
	runErr := ctx.Run()

	// Close writers to finish copies
	_ = wOut.Close()
	_ = wErr.Close()
	<-doneOut
	<-doneErr

	return bufOut.String(), bufErr.String(), runErr
}

func TestFactsTierAndNetWithNames(t *testing.T) {
	out, errOut, err := runCLI([]string{"facts", "-t", "2", "--net", "os.name", "os.version"})
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if want := "facts: tier=2 net=true (stub)"; !strings.Contains(out, want) {
		t.Fatalf("stdout missing %q, got: %q", want, out)
	}
	// Should also include the provided positional fact names.
	if !strings.Contains(out, "os.name") || !strings.Contains(out, "os.version") {
		t.Fatalf("stdout missing provided fact names, got: %q", out)
	}
	if errOut != "" {
		t.Fatalf("unexpected stderr: %q", errOut)
	}
}

func TestCollectorStub(t *testing.T) {
	out, errOut, err := runCLI([]string{"collector"})
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if want := "collector: running (stub)"; !strings.Contains(out, want) {
		t.Fatalf("stdout missing %q, got: %q", want, out)
	}
	if errOut != "" {
		t.Fatalf("unexpected stderr: %q", errOut)
	}
}

func TestRequiresStub(t *testing.T) {
	out, errOut, err := runCLI([]string{"requires", "list"})
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if want := "requires: listing (stub)"; !strings.Contains(out, want) {
		t.Fatalf("stdout missing %q, got: %q", want, out)
	}
	if errOut != "" {
		t.Fatalf("unexpected stderr: %q", errOut)
	}
}

func TestRequiresCheck(t *testing.T) {
	out, errOut, err := runCLI([]string{"requires", "check"})
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if want := "requires: checking (stub)"; !strings.Contains(out, want) {
		t.Fatalf("stdout missing %q, got: %q", want, out)
	}
	if errOut != "" {
		t.Fatalf("unexpected stderr: %q", errOut)
	}
}

func TestDebugFlag(t *testing.T) {
	out, errOut, err := runCLI([]string{"--debug", "facts", "-t", "1"})
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if want := "debug: enabled"; !strings.Contains(errOut, want) {
		t.Fatalf("stderr missing %q, got: %q", want, errOut)
	}
	if want := "facts: tier=1 net=false (stub)"; !strings.Contains(out, want) {
		t.Fatalf("stdout missing %q, got: %q", want, out)
	}
}

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	cli "github.com/urfave/cli/v3"
)

// buildCmd replicates the CLI from main() for testing without os.Exit.
func buildCmd() *cli.Command {
	return &cli.Command{
		Name:  "clu",
		Usage: "urfave/cli v3 example with facts/collector/requires",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "debug", Usage: "Enable debug logging."},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			if cmd.Bool("debug") {
				fmt.Fprintln(os.Stderr, "debug: enabled")
			}
			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name:  "facts",
				Usage: "Show facts (stub)",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "tier", Aliases: []string{"t"}, Usage: "Tier level (1, 2, or 3).", Value: 1},
					&cli.BoolFlag{Name: "net", Usage: "Enable network access."},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					tier := cmd.Int("tier")
					if tier < 1 || tier > 3 {
						return fmt.Errorf("invalid --tier %d (allowed: 1, 2, 3)", tier)
					}
					names := cmd.Args().Slice()
					if len(names) > 0 {
						fmt.Printf("facts: tier=%d net=%v (stub) %s\n", tier, cmd.Bool("net"), strings.Join(names, " "))
					} else {
						fmt.Printf("facts: tier=%d net=%v (stub)\n", tier, cmd.Bool("net"))
					}
					return nil
				},
			},
			{
				Name:  "collector",
				Usage: "Run collector (stub)",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("collector: running (stub)")
					return nil
				},
			},
			{
				Name:  "requires",
				Usage: "Requires actions: list or check",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					args := cmd.Args().Slice()
					if len(args) < 1 {
						return fmt.Errorf("requires expects one positional argument: list or check")
					}
					mode := args[0]
					switch mode {
					case "list":
						fmt.Println("requires: listing (stub)")
					case "check":
						fmt.Println("requires: checking (stub)")
					default:
						return fmt.Errorf("invalid requires mode %q (allowed: list, check)", mode)
					}
					return nil
				},
			},
		},
	}
}

func runCLI(args []string) (stdout, stderr string, err error) {
	cmd := buildCmd()
	args = append([]string{"clu"}, args...)

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

	runErr := cmd.Run(context.Background(), args)

	_ = wOut.Close()
	_ = wErr.Close()
	<-doneOut
	<-doneErr

	return bufOut.String(), bufErr.String(), runErr
}

func TestFactsTierAndNetWithNames(t *testing.T) {
	out, errOut, err := runCLI([]string{"facts", "--tier", "2", "--net", "os.name", "os.version"})
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if want := "facts: tier=2 net=true (stub)"; !strings.Contains(out, want) {
		t.Fatalf("stdout missing %q, got: %q", want, out)
	}
	if !strings.Contains(out, "os.name") || !strings.Contains(out, "os.version") {
		t.Fatalf("stdout missing provided fact names, got: %q", out)
	}
	if errOut != "" {
		t.Fatalf("unexpected stderr: %q", errOut)
	}
}

func TestFactsTierShortAlias(t *testing.T) {
	out, errOut, err := runCLI([]string{"facts", "-t", "3"})
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if want := "facts: tier=3 net=false (stub)"; !strings.Contains(out, want) {
		t.Fatalf("stdout missing %q, got: %q", want, out)
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

func TestRequiresStubList(t *testing.T) {
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

func TestRequiresStubCheck(t *testing.T) {
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
	out, errOut, err := runCLI([]string{"--debug", "facts", "--tier", "1"})
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

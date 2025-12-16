package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	cli "github.com/urfave/cli"
)

// buildApp replicates the CLI from main() for testing without os.Exit.
func buildApp() *cli.App {
	app := cli.NewApp()
	app.Name = "clu"
	app.Usage = "urfave/cli v1 example with facts/collector/requires"
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "debug", Usage: "Enable debug logging."},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			fmt.Fprintln(os.Stderr, "debug: enabled")
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "facts",
			Usage: "Show facts (stub)",
			Flags: []cli.Flag{
				cli.IntFlag{Name: "tier, t", Usage: "Tier level (1, 2, or 3).", Value: 1},
				cli.StringFlag{Name: "out", Usage: "Output format: dots, json, or shell.", Value: "dots"},
				cli.BoolFlag{Name: "net", Usage: "Enable network access."},
			},
			Action: func(c *cli.Context) error {
				tier := c.Int("tier")
				if tier < 1 || tier > 3 {
					return fmt.Errorf("invalid --tier %d (allowed: 1, 2, 3)", tier)
				}
				outfmt := c.String("out")
				switch outfmt {
				case "dots", "json", "shell":
				default:
					return fmt.Errorf("invalid --out %q (allowed: dots, json, shell)", outfmt)
				}
				names := c.Args()
				if len(names) > 0 {
					fmt.Printf("facts: tier=%d net=%v (stub) %s\n", tier, c.Bool("net"), strings.Join(names, " "))
				} else {
					fmt.Printf("facts: tier=%d net=%v (stub)\n", tier, c.Bool("net"))
				}
				return nil
			},
		},
		{
			Name:  "collector",
			Usage: "Run collector (stub)",
			Action: func(c *cli.Context) error {
				fmt.Println("collector: running (stub)")
				return nil
			},
		},
		{
			Name:  "requires",
			Usage: "Requires actions: list or check",
			Action: func(c *cli.Context) error {
				args := c.Args()
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
	}
	return app
}

func runCLI(args []string) (stdout, stderr string, err error) {
	app := buildApp()

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

	err = app.Run(append([]string{"clu"}, args...))

	_ = wOut.Close()
	_ = wErr.Close()
	<-doneOut
	<-doneErr

	return bufOut.String(), bufErr.String(), err
}

func TestFactsTierAndNetWithNames(t *testing.T) {
	out, errOut, err := runCLI([]string{"facts", "-t", "2", "--net", "os.name", "os.version"})
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

func TestRequiresList(t *testing.T) {
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

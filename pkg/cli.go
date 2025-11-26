package pkg

// Go port of src/clu/cli.py using pflag for argument parsing.
// This provides a main-style entrypoint plus subcommand registration similar
// to the Python implementation. Subcommand implementations (report, archive,
// requires) are expected to be added separately.

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/pflag"
)

// Version variables could be set at build time via -ldflags.
var (
	MinimumGoVersion = "1.22" // analogous to Python minimum in original
)

// Globals capturing parsed state (similar to Python set_config)
var (
	DebugLevel   int
	VerboseLevel int
)

// CommandFunc is a function run for a subcommand.
type CommandFunc func() int

// registry of subcommands
var commands = map[string]CommandFunc{}

// RegisterCommand adds a subcommand by name.
func RegisterCommand(name string, fn CommandFunc) {
	commands[name] = fn
}

// setupLogging is a placeholder to integrate with a logging framework.
// For now it does nothing beyond being a stub.
func setupLogging(debug int) {
	// Integrate with chosen logging library later.
}

// parseCmdline configures and parses flags and returns the selected command function.
func parseCmdline(args []string) CommandFunc {
	fs := pflag.NewFlagSet("clu", pflag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: clu [options] <command>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nCommands:\n")
		for name := range commands {
			fmt.Fprintf(os.Stderr, "  %s\n", name)
		}
	}

	fs.CountVarP(&DebugLevel, "debug", "d", "Increase debugging output (can be used twice).")
	fs.CountVarP(&VerboseLevel, "verbose", "v", "Increase verbosity of output.")

	showVersion := fs.Bool("version", false, "Show version and exit")

	// Parse flags up to first non-flag (the command)
	fs.Parse(args)

	if *showVersion {
		fmt.Printf("clu %s\n", Version)
		os.Exit(0)
	}

	remaining := fs.Args()
	cmdName := "report" // default command analogous to Python
	if len(remaining) > 0 {
		cmdName = remaining[0]
		remaining = remaining[1:]
	}

	cmdFn, ok := commands[cmdName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmdName)
		fs.Usage()
		os.Exit(2)
	}

	// TODO: Pass remaining args to subcommand-specific parsing once implemented.
	return cmdFn
}

// Main is the entry point analogous to Python main().
func Main() int {
	// Basic runtime version check placeholder (Go rarely needs this the same way).
	_ = runtime.Version()

	cmdFn := parseCmdline(os.Args[1:])
	setupLogging(DebugLevel)

	// In Python: log start & debug command line; we omit until logging integrated.
	return cmdFn()
}

// If building a standalone binary, package main could call pkg.Main().

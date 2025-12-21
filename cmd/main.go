package main

// Go port of src/clu/cli.py using pflag for argument parsing.
// This provides a main-style entrypoint plus subcommand registration similar
// to the Python implementation. Subcommand implementations (report, archive,
// requires) are expected to be added separately.

import (
	"flag"
	"fmt"
	"os"

	"github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/subcmd"
)

// setupLogging is a placeholder to integrate with a logging framework.
// For now it does nothing beyond being a stub.
// func setupLogging() {
// 	fmt.Printf("LOG: setupLogging called with debug=%t\n", pkg.Config.Debug)
// 	fmt.Printf("LOG: setupLogging called with verbose=%t\n", pkg.Config.Verbose)
// }

// parseCmdline configures and parses flags and returns the selected command function.
func parseCmdline(args []string) {
	fmt.Printf("LOG: parseCmdline called with args=%v\n", args)

	flag.BoolVar(&pkg.Config.Debug, "debug", false, "enable debug output")
	flag.BoolVar(&pkg.Config.Verbose, "verbose", false, "enable verbose output")

	flag.Parse()
}

func main() {

	parseCmdline(os.Args[1:])
	// setupLogging()

	exit := subcmd.ReportFacts()

	os.Exit(exit)
}

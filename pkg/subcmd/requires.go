package subcmd

// Go port of src/clu/cmd/requires.py excluding parse_args. Provides listing and
// checking of system requirements gathered from the OpSys abstraction. Caller
// is expected to decide which subcommand (list or check) to invoke.

import (
	"fmt"
	"os"
	"strings"

	"github.com/NHGRI/clu/pkg/facts"
	"github.com/NHGRI/clu/pkg/input"
)

// RequiresCmd implements the "requires" subcommand (stub only).
type RequiresCmd struct {
	Mode string `arg:"" enum:"list,check" help:"Operation to perform: list or check."`
}

func (r *RequiresCmd) Run(stdout input.Stdout, stderr input.Stderr) error {
	switch r.Mode {
	case "list":
		// fmt.Println("requires: listing (stub)")
		listRequires(stdout, stderr)
	case "check":
		// fmt.Println("requires: checking (stub)")
		checkRequires(stdout, stderr)
	}
	return nil
}

// listRequires lists all file and program requirements.
func listRequires(stdout input.Stdout, stderr input.Stderr) int {
	reqs := facts.OpSysFactory().Requires()
	fmt.Fprintln(stdout, "Listing Requirements:")
	fmt.Fprintln(stdout, "----------------------")
	fmt.Fprintln(stdout, "Files:")
	for _, file := range reqs.Files {
		fmt.Fprintf(stdout, "  - %s\n", file)
	}
	fmt.Fprintln(stdout, "Programs:")
	for _, prog := range reqs.Programs {
		fmt.Fprintf(stdout, "  - %s\n", prog)
	}
	// APIs omitted; not used currently.
	return 0
}

// checkRequires checks existence of required files and programs.
func checkRequires(stdout input.Stdout, stderr input.Stderr) int {
	reqs := facts.OpSysFactory().Requires()
	fmt.Fprintln(stdout, "Checking requirements:")
	fmt.Fprintln(stdout, "----------------------")
	fmt.Fprintln(stdout, "Files:")
	for _, file := range reqs.Files {
		if checkFileExists(file) {
			fmt.Fprintf(stdout, "  - [ok] %s\n", file)
		} else {
			fmt.Fprintf(stdout, "  - [MISSING] %s\n", file)
		}
	}
	fmt.Fprintln(stdout, "Programs:")
	for _, prog := range reqs.Programs {
		if checkProgramExists(prog) {
			fmt.Fprintf(stdout, "  - [ok] %s\n", prog)
		} else {
			fmt.Fprintf(stdout, "  - [MISSING] %s\n", prog)
		}
	}
	return 0
}

// checkFileExists mimics Python check_file_exists helper using os.Stat.
func checkFileExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

// checkProgramExists mimics Python check_program_exists by searching PATH.
func checkProgramExists(name string) bool {
	if name == "" {
		return false
	}
	// If the requirement includes spaces (arguments), use first token for lookup.
	prog := name
	for i, r := range name {
		if r == ' ' {
			prog = name[:i]
			break
		}
	}
	// Look through PATH entries.
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return false
	}
	for _, dir := range splitPath(pathEnv) {
		candidate := dir + string(os.PathSeparator) + prog
		if fi, err := os.Stat(candidate); err == nil && !fi.IsDir() {
			return true
		}
	}
	return false
}

// splitPath splits PATH variable; separated for testability.
func splitPath(p string) []string { return strings.Split(p, string(os.PathListSeparator)) }

package cmd

// Go port of src/clu/cmd/requires.py excluding parse_args. Provides listing and
// checking of system requirements gathered from the OpSys abstraction. Caller
// is expected to decide which subcommand (list or check) to invoke.

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/huntermatthews/clu/pkg/opsys"
)

// RequiresConfig carries the chosen subcommand name ("list" or "check").
type RequiresConfig struct {
	Subcmd string
}

// Run dispatches to list or check logic based on Subcmd. Returns exit code.
func (c *RequiresConfig) Run() int {
	switch c.Subcmd {
	case "list":
		return listRequires()
	case "check":
		return checkRequires()
	default:
		fmt.Fprintf(os.Stderr, "Unknown sub-command: %s\n", c.Subcmd)
		return 2
	}
}

// opsysFactory replicates Python opsys_factory minimal logic using runtime.GOOS.
func opsysFactory() *opsys.OpSys {
	if runtime.GOOS == "darwin" {
		return opsys.NewDarwin()
	}
	return opsys.NewLinux()
}

// listRequires lists all file and program requirements.
func listRequires() int {
	reqs := opsysFactory().Requires()
	fmt.Println("Listing Requirements:")
	fmt.Println("----------------------")
	fmt.Println("Files:")
	for _, file := range reqs.Files {
		fmt.Printf("  - %s\n", file)
	}
	fmt.Println("Programs:")
	for _, prog := range reqs.Programs {
		fmt.Printf("  - %s\n", prog)
	}
	// APIs omitted; not used currently.
	return 0
}

// checkRequires checks existence of required files and programs.
func checkRequires() int {
	reqs := opsysFactory().Requires()
	fmt.Println("Checking requirements:")
	fmt.Println("----------------------")
	fmt.Println("Files:")
	for _, file := range reqs.Files {
		if checkFileExists(file) {
			fmt.Printf("  - [ok] %s\n", file)
		} else {
			fmt.Printf("  - [MISSING] %s\n", file)
		}
	}
	fmt.Println("Programs:")
	for _, prog := range reqs.Programs {
		if checkProgramExists(prog) {
			fmt.Printf("  - [ok] %s\n", prog)
		} else {
			fmt.Printf("  - [MISSING] %s\n", prog)
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

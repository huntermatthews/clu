package pkg

// Port of src/clu/input.py providing utilities for reading text files and running
// simple text-based programs, plus helper transformations.

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const FileSizeLimit = 1 * 1024 * 1024 // 1MB arbitrary limit

// TextFile reads a file and returns its contents as a string. If optional is true
// and the file does not exist, an empty string is returned without treating it as an error.
func TextFile(fname string, optional bool) string {
	info, err := os.Stat(fname)
	if err != nil {
		if optional {
			log.Printf("optional file not found: %s", fname)
		} else {
			log.Printf("file not found: %s", fname)
		}
		return ""
	}
	if info.IsDir() {
		log.Printf("path is a directory, not file: %s", fname)
		return ""
	}
	if info.Size() > FileSizeLimit {
		log.Printf("file size exceeds limit: %d > %d", info.Size(), FileSizeLimit)
		return ""
	}
	data, err := os.ReadFile(fname)
	if err != nil {
		log.Printf("error reading file %s: %v", fname, err)
		return ""
	}
	return string(data)
}

// TextProgram runs a command line (split on spaces) and returns stdout and the exit code.
// On   error returns empty output and rc=1.
func TextProgram(cmdline string) (string, int) {
	parts := strings.Fields(cmdline)
	if len(parts) == 0 {
		return "", 1
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	out, err := cmd.Output()
	if err != nil {
		// Try to fetch combined output for richer debugging
		if ee, ok := err.(*exec.ExitError); ok {
			return string(ee.Stderr), ee.ExitCode()
		}
		log.Printf("error running program %s: %v", cmdline, err)
		return "", 1
	}
	return string(out), 0
}

// TransformCmdlineToFilename converts a command line into a safe-ish filename.
// Spaces -> underscores, slashes -> percent signs. Returns the filename and
// the filename with _rc suffix (akin to Python returning tuple of name and rc name).
func TransformCmdlineToFilename(cmdline string) (string, string) {
	name := strings.ReplaceAll(cmdline, " ", "_")
	name = strings.ReplaceAll(name, "/", "%")
	return name, name + "_rc"
}

// CheckProgramExists returns the absolute path to the program if found in PATH,
// otherwise empty string.
func CheckProgramExists(program string) string {
	first := strings.Fields(program)
	if len(first) == 0 {
		return ""
	}
	path, err := exec.LookPath(first[0])
	if err != nil {
		return ""
	}
	return path
}

// CheckFileExists returns the absolute path if the file exists, else empty string.
func CheckFileExists(fname string) string {
	abs, err := filepath.Abs(fname)
	if err != nil {
		return ""
	}
	info, err := os.Stat(abs)
	if err != nil || info.IsDir() {
		return ""
	}
	return abs
}

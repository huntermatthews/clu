package input

// Port of src/clu/input.py providing utilities for reading text files and running
// simple text-based programs, plus helper transformations.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const OutputSizeLimit = 1 * 1024 * 1024 // 1MB arbitrary limit

// TextProgram runs a command line (split on spaces) and returns combined
// stdout/stderr output, the exit code, and an error. If the process runs and
// exits with a non-zero status the returned error will be of type
// *exec.ExitError (so callers can inspect it), and the exit code will be
// returned as the second value.
func TextProgram(cmdline string) (string, int, error) {
	parts := strings.Fields(cmdline)
	if len(parts) == 0 {
		return "", 1, fmt.Errorf("empty command line")
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	out, err := cmd.CombinedOutput()
	output := string(out)

	// Derive an exit code we can return even if the process failed.
	exitCode := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else {
			exitCode = 1
		}
	}

	// Enforce a maximum output size similar to TextFile's FileSizeLimit.
	// BUG: we need to stop at the limit while reading output, not after.
	if len(out) > OutputSizeLimit {
		return output, exitCode, fmt.Errorf("output size exceeds limit: %d > %d", len(out), OutputSizeLimit)
	}

	return output, 0, nil
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

// TextFile reads a text file and returns its contents as a string. If the file
// does not exist or cannot be read, an error is returned. If the file is empty
// an error is also returned. Files larger than OutputSizeLimit are rejected.
func TextFile(fname string) (string, error) {
	info, err := os.Stat(fname)
	if err != nil {
		if os.IsNotExist(err) {
			return "", os.ErrNotExist
		}
		return "", fmt.Errorf("stat %s: %w", fname, err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("path is a directory: %s", fname)
	}

	if info.Size() > OutputSizeLimit {
		return "", fmt.Errorf("file size exceeds limit: %d > %d", info.Size(), OutputSizeLimit)
	}

	data, err := os.ReadFile(fname)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", fname, err)
	}

	// FIXME: consider whether empty files are truly an error in all contexts
	// FIXME: base this off the stat call instead of reading the file?
	if len(data) == 0 {
		return "", fmt.Errorf("empty file content")
	}
	return string(data), nil
}

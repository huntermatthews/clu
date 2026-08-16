// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package input

// MockTextFile reads a file from a per-host mock directory set in global.Config.MockDir.
// The mock directory is expected to contain files laid out by their absolute
// path (e.g. MockDir + "/proc/cpuinfo"). The function mirrors the behavior of
// `TextFile` but does not call it — it implements the same error checks and
// size limits and returns the same error semantics.

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/huntermatthews/clu/internal/global"
)

// MockTextFile reads the mock file corresponding to `fname` under the
// configured `global.Config.MockDir`. `fname` is treated like an absolute
// path; its leading path separator will be removed before joining with the
// mock directory to ensure the mock file layout mirrors the real filesystem.
func MockTextFile(fname string) (string, error) {
	// Build the path under the mock dir. If MockDir is empty, fall back to fname.
	if global.Config.MockDir == "" {
		panic("global.Config.MockDir is not set for MockTextFile")
	}
	cleaned := strings.TrimPrefix(fname, string(os.PathSeparator))
	path := filepath.Join(global.Config.MockDir, cleaned)

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", os.ErrNotExist
		}
		return "", fmt.Errorf("stat %s: %w", path, err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("path is a directory: %s", path)
	}

	if info.Size() > OutputSizeLimit {
		return "", fmt.Errorf("file size exceeds limit: %d > %d", info.Size(), OutputSizeLimit)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}

	if len(data) == 0 {
		return "", fmt.Errorf("empty file content")
	}
	return string(data), nil
}

// MockTextProgram returns mocked program output and return code for a given
// command line. It looks up the mock files under
// `global.Config.MockDir/_programs/<name>` where <name> comes from
// `TransformCmdlineToFilename(cmdline)`; the output file is `<name>` and the
// return-code file is `<name>_rc`. Both files are optional; if the output file
// is missing, this function returns an error similar to running the command
// failing to produce output. If the rc file is missing, rc defaults to 0.
func MockTextProgram(cmdline string) (string, int, error) {
	if global.Config.MockDir == "" {
		panic("global.Config.MockDir is not set for MockTextProgram")
	}
	if strings.TrimSpace(cmdline) == "" {
		return "", 1, fmt.Errorf("empty command line")
	}

	name, rcName := TransformCmdlineToFilename(cmdline)
	// program outputs are under _programs directory
	base := filepath.Join(global.Config.MockDir, "_programs")
	outPath := filepath.Join(base, name)
	rcPath := filepath.Join(base, rcName)

	// Read output file
	outInfo, outErr := os.Stat(outPath)
	if outErr != nil {
		if os.IsNotExist(outErr) {
			return "", 1, os.ErrNotExist
		}
		return "", 1, fmt.Errorf("stat %s: %w", outPath, outErr)
	}
	if outInfo.IsDir() {
		return "", 1, fmt.Errorf("path is a directory: %s", outPath)
	}
	if outInfo.Size() > OutputSizeLimit {
		return "", 1, fmt.Errorf("output size exceeds limit: %d > %d", outInfo.Size(), OutputSizeLimit)
	}
	outData, err := os.ReadFile(outPath)
	if err != nil {
		return "", 1, fmt.Errorf("read %s: %w", outPath, err)
	}

	// Default rc
	rc := 0
	if rcInfo, err := os.Stat(rcPath); err == nil && !rcInfo.IsDir() {
		// read rc file; if reading fails, propagate error
		rcData, err := os.ReadFile(rcPath)
		if err != nil {
			return "", 1, fmt.Errorf("read %s: %w", rcPath, err)
		}
		s := strings.TrimSpace(string(rcData))
		if s != "" {
			v, perr := strconv.Atoi(s)
			if perr != nil {
				panic(fmt.Errorf("invalid rc file %s: %q", rcPath, s))
			}
			rc = v
		}
	}

	return string(outData), rc, nil
}

// MockProgramCheck checks if a program exists in the mock environment by
// looking for the program's output file under
// `global.Config.MockDir/_programs/<name>` where <name> comes from
// `TransformCmdlineToFilename(cmdline)`. Returns the program name (first
// component before underscore) if found, otherwise empty string.
func MockProgramCheck(cmdline string) string {
	if global.Config.MockDir == "" {
		panic("global.Config.MockDir is not set for MockProgramCheck")
	}
	if strings.TrimSpace(cmdline) == "" {
		return ""
	}

	name, _ := TransformCmdlineToFilename(cmdline)
	// program outputs are under _programs directory
	base := filepath.Join(global.Config.MockDir, "_programs")
	outPath := filepath.Join(base, name)

	// Check if output file exists and is a regular file
	info, err := os.Stat(outPath)
	if err != nil {
		return ""
	}
	if info.IsDir() {
		return ""
	}

	// Return only the program name (first part before underscore)
	progName := strings.SplitN(name, "_", 2)[0]
	return progName
}

// MockTextFileAge returns the mocked modification time for a file by reading
// an RFC3339 timestamp from a file named `<fname>_mtime` under the configured
// `global.Config.MockDir`. The mtime file should contain an RFC3339 timestamp
// string (e.g., "2026-02-04T15:30:00Z"). If the mtime file is missing, it defaults
// to the actual modification time of the mock file itself.
func MockTextFileAge(fname string) (time.Time, error) {
	if global.Config.MockDir == "" {
		panic("global.Config.MockDir is not set for MockTextFileAge")
	}

	// Build the path under the mock dir
	cleaned := strings.TrimPrefix(fname, string(os.PathSeparator))
	path := filepath.Join(global.Config.MockDir, cleaned)
	mtimePath := path + "_mtime"

	// Check that the base file exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return time.Time{}, os.ErrNotExist
		}
		return time.Time{}, fmt.Errorf("stat %s: %w", path, err)
	}
	if info.IsDir() {
		return time.Time{}, fmt.Errorf("path is a directory: %s", path)
	}

	// Try to read the mtime file
	if mtimeInfo, err := os.Stat(mtimePath); err == nil && !mtimeInfo.IsDir() {
		mtimeData, err := os.ReadFile(mtimePath)
		if err != nil {
			return time.Time{}, fmt.Errorf("read %s: %w", mtimePath, err)
		}
		s := strings.TrimSpace(string(mtimeData))
		if s != "" {
			timestamp, perr := time.Parse(time.RFC3339, s)
			if perr != nil {
				return time.Time{}, fmt.Errorf("invalid mtime file %s: %q", mtimePath, s)
			}
			return timestamp, nil
		}
	}

	// Default to the actual file's modification time
	return info.ModTime(), nil
}

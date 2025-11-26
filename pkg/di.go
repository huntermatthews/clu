package pkg

import "fmt"

// Dependency injection points for external command and file access to enable deterministic tests.
// Tests can override these variables to supply canned outputs without invoking the real system.

// CommandRunnerFunc executes a command line and returns stdout + exit code.
type CommandRunnerFunc func(cmdline string) (string, int)

// FileReaderFunc reads a file path returning contents and error.
type FileReaderFunc func(path string) (string, error)

// CommandRunner defaults to TextProgram; override in tests as needed.
var CommandRunner CommandRunnerFunc = func(cmdline string) (string, int) { return TextProgram(cmdline) }

// FileReader defaults to using TextFile (optional=false); override in tests for custom contents.
var FileReader FileReaderFunc = func(path string) (string, error) {
	data := TextFile(path, false)
	if data == "" { // treat empty as error for deterministic signaling
		return "", ErrEmptyFile
	}
	return data, nil
}

// ErrEmptyFile returned when FileReader encounters an empty file using the default implementation.
var ErrEmptyFile = fmt.Errorf("empty file content")

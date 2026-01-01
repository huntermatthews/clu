package input

// Dependency injection points for external command and file access to enable deterministic tests.
// Tests can override these variables to supply canned outputs without invoking the real system.

import "io"

// CommandRunnerFunc executes a command line and returns stdout, exit code and error.
type CommandRunnerFunc func(cmdline string) (string, int, error)

// FileReaderFunc reads a file path returning contents and error.
type FileReaderFunc func(path string) (string, error)

// CommandRunner defaults to TextProgram; override in tests as needed.
var CommandRunner CommandRunnerFunc = TextProgram

// FileReader defaults to using TextFile; override in tests for custom contents.
var FileReader FileReaderFunc = TextFile

// Stdout is a marker type for dependency injection of stdout.
// We use a distinct type because Kong resolves dependencies by type; binding io.Writer
// twice would cause the second binding to overwrite the first.
type Stdout io.Writer

// Stderr is a marker type for dependency injection of stderr.
// We use a distinct type because Kong resolves dependencies by type; binding io.Writer
// twice would cause the second binding to overwrite the first.
type Stderr io.Writer

package pkg

// Dependency injection points for external command and file access to enable deterministic tests.
// Tests can override these variables to supply canned outputs without invoking the real system.

// CommandRunnerFunc executes a command line and returns stdout, exit code and error.
type CommandRunnerFunc func(cmdline string) (string, int, error)

// FileReaderFunc reads a file path returning contents and error.
type FileReaderFunc func(path string) (string, error)

// CommandRunner defaults to TextProgram; override in tests as needed.
var CommandRunner CommandRunnerFunc = TextProgram

// FileReader defaults to using TextFile; override in tests for custom contents.
var FileReader FileReaderFunc = TextFile

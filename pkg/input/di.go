// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package input

// Dependency injection points for external command and file access to enable deterministic tests.
// Tests can override these variables to supply canned outputs without invoking the real system.

import (
	"io"
	"time"
)

// CommandRunnerFunc executes a command line and returns stdout, exit code and error.
type CommandRunnerFunc func(cmdline string) (string, int, error)

// FileReaderFunc reads a file path returning contents and error.
type FileReaderFunc func(path string) (string, error)

// FileAgeReaderFunc reads a file's modification time returning time.Time and error.
type FileAgeReaderFunc func(path string) (time.Time, error)

type ProgramCheckerFunc func(cmdline string) string

// CommandRunner defaults to TextProgram; override in tests as needed.
var CommandRunner CommandRunnerFunc = TextProgram

// FileReader defaults to using TextFile; override in tests for custom contents.
var FileReader FileReaderFunc = TextFile

// FileAgeReader defaults to using TextFileAge; override in tests for custom times.
var FileAgeReader FileAgeReaderFunc = TextFileAge

var ProgramChecker ProgramCheckerFunc = CheckProgramExists

// Stdout is a marker type for dependency injection of stdout.
// We use a distinct type because Kong resolves dependencies by type; binding io.Writer
// twice would cause the second binding to overwrite the first.
type Stdout io.Writer

// Stderr is a marker type for dependency injection of stderr.
// We use a distinct type because Kong resolves dependencies by type; binding io.Writer
// twice would cause the second binding to overwrite the first.
type Stderr io.Writer

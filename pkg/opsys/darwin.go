package opsys

// Go port of src/clu/opsys/darwin.py providing the macOS (Darwin) operating system
// source aggregation and default/early fact lists.

import (
	sources "github.com/huntermatthews/clu/pkg/source"
)

// NewDarwin constructs the Darwin OpSys with its ordered sources and fact lists.
// Order matches the Python implementation for deterministic precedence.
func NewDarwin() *OpSys {
	srcs := []sources.Source{
		&sources.Uname{},
		&sources.SystemVersionPlist{},
		&sources.MacOSName{},
		&sources.Uptime{},
		&sources.Clu{},
	}
	defaults := []string{"os.name", "os.hostname", "os.version", "os.code_name", "run.uptime", "clu.version"}
	early := []string{"os.version"}
	return New(srcs, defaults, early)
}

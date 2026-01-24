package facts

// Go port of src/clu/opsys/darwin.py providing the macOS (Darwin) operating system
// source aggregation and early fact lists.

import (
	"github.com/NHGRI/clu/pkg/facts/sources"
	"github.com/NHGRI/clu/pkg/facts/types"
)

// NewDarwin constructs the Darwin OpSys with its ordered sources and fact lists.
// Order matches the Python implementation for deterministic precedence.
func NewDarwin() *OpSys {
	return &OpSys{
		Sources: []types.Sources{
			&sources.Uname{},
			&sources.SystemVersionPlist{},
			&sources.MacOSName{},
			&sources.Uptime{},
			&sources.Clu{},
		},
		EarlyFacts: []string{"os.version"},
	}
}

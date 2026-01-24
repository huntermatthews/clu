package facts

// Windows operating system source aggregation and early fact lists.

import (
	"github.com/NHGRI/clu/pkg/facts/sources"
	"github.com/NHGRI/clu/pkg/facts/types"
)

// NewWindows constructs the Windows OpSys with its ordered sources and fact lists.
func NewWindows() *OpSys {
	return &OpSys{
		Sources: []types.Sources{
			&sources.WindowsSysteminfo{},
			&sources.Clu{},
		},
		EarlyFacts: []string{},
	}
}

package facts

// Windows operating system source aggregation and default/early fact lists.

import (
	"github.com/NHGRI/clu/pkg/facts/sources"
	"github.com/NHGRI/clu/pkg/facts/types"
)

// NewWindows constructs the Windows OpSys with its ordered sources and fact lists.
func NewWindows() *OpSys {
	srcs := []types.Sources{
		//		&sources.WindowsSystemInfo{},
		&sources.WindowsVer{},
		&sources.Clu{},
	}
	defaults := []string{
		"os.name",
		"os.version",
		"clu.version",
	}
	early := []string{}
	return NewOpSys(srcs, defaults, early)
}

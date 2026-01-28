package facts

// Constructs the Linux OpSys with its ordered fact sources and early fact lists.

import (
	"github.com/NHGRI/clu/pkg/facts/sources"
	"github.com/NHGRI/clu/pkg/facts/types"
)

// NewLinux constructs the Linux OpSys. Order mirrors Python minus unported AwsImds source.
func NewLinux() *OpSys {
	return &OpSys{
		Sources: []types.Sources{
			&sources.Clu{},
			&sources.LinuxCheckUpdate{},
			&sources.IpAddr{},
			&sources.Ipmitool{},
			&sources.Lscpu{},
			&sources.Lsmem{},
			&sources.NoSalt{},
			&sources.OsRelease{},
			&sources.ProcCpuinfo{},
			&sources.ProcUptime{},
			&sources.Selinux{},
			&sources.SysDmi{},
			&sources.Uname{},
			&sources.VirtWhat{},
			//		&sources.UdevadmRam{},
		},
		EarlyFacts: []string{
			"phy.arch",
			"phy.platform",
		},
	}
}

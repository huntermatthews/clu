package facts

// Constructs the Linux OpSys with its ordered fact sources and default/early fact lists.

import (
	"github.com/NHGRI/clu/pkg/facts/sources"
	"github.com/NHGRI/clu/pkg/facts/types"
)

// NewLinux constructs the Linux OpSys. Order mirrors Python minus unported AwsImds source.
func NewLinux() *OpSys {
	srcs := []types.Sources{
		&sources.Clu{},
		&sources.DnfCheckUpdate{},
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
	}
	defaults := []string{
		"os.name",
		"os.hostname",
		"os.distro.name",
		"os.distro.version",
		"phy.platform",
		"phy.ram",
		"run.uptime",
		"clu.version",
	}
	early := []string{
		"phy.arch",
		"phy.platform",
	}
	return NewOpSys(srcs, defaults, early)
}

package opsys

// Go port of src/clu/opsys/linux.py
// Constructs the Linux OpSys with its ordered fact sources and default/early fact lists.
// Note: Source AwsImds not yet ported.

import (
	"github.com/huntermatthews/clu/pkg/sources"
)

// NewLinux constructs the Linux OpSys. Order mirrors Python minus unported AwsImds source.
func NewLinux() *OpSys {
	srcs := []sources.Source{
		// &sources.AwsImds{}, // TODO: port aws_imds.py
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
		&sources.UdevadmRam{},
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
	early := []string{"phy.arch", "phy.platform"}
	return New(srcs, defaults, early)
}

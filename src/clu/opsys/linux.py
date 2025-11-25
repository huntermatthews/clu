import logging

from clu.opsys import OpSys
from clu.sources import (
    aws_imds,
    clu,
    dnf_checkupdate,
    ip_addr,
    ipmitool,
    lscpu,
    lsmem,
    no_salt,
    os_release,
    proc_cpuinfo,
    proc_uptime,
    selinux,
    sys_dmi,
    uname,
    virt_what,
)

log = logging.getLogger(__name__)


class Linux(OpSys):
    _sources = [
        aws_imds.AwsImds(),
        clu.Clu(),
        dnf_checkupdate.DnfCheckUpdate(),
        ip_addr.IpAddr(),
        ipmitool.Ipmitool(),
        lscpu.Lscpu(),
        lsmem.Lsmem(),
        no_salt.NoSalt(),
        os_release.OsRelease(),
        proc_cpuinfo.ProcCpuinfo(),
        proc_uptime.ProcUptime(),
        selinux.Selinux(),
        sys_dmi.SysDmi(),
        uname.Uname(),
        virt_what.VirtWhat(),
    ]

    def default_facts(self) -> list:
        """Default facts for Linux."""
        return [
            "os.name",
            "os.hostname",
            "os.distro.name",
            "os.distro.version",
            "phy.platform",
            "phy.ram",
            "run.uptime",
            "clu.version",
        ]

    def early_facts(self) -> list:
        return ["phy.arch", "phy.platform"]

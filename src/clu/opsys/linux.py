import logging

from clu.opsys import OpSys
from clu.sources import (
    clu,
    dnf_checkupdate,
    ip_addr,
    ipmitool,
    lscpu,
    no_salt,
    os_release,
    proc_cpuinfo,
    proc_uptime,
    selinux,
    sys_dmi,
    udevadm_ram,
    uname,
    virt_what,
)

log = logging.getLogger(__name__)


class Linux(OpSys):
    _sources = [
        clu.Clu(),
        dnf_checkupdate.DnfCheckUpdate(),
        ip_addr.IpAddr(),
        ipmitool.Ipmitool(),
        lscpu.Lscpu(),
        no_salt.NoSalt(),
        os_release.OsRelease(),
        proc_cpuinfo.ProcCpuinfo(),
        proc_uptime.ProcUptime(),
        selinux.Selinux(),
        sys_dmi.SysDmi(),
        udevadm_ram.UdevadmRam(),
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
            "phy.ram",
            "run.uptime",
            "clu.version",
        ]

    def early_facts(self) -> list:
        return ["phy.arch", "phy.platform"]

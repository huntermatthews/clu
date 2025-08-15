import logging

from clu.opsys import OpSys
from clu import Facts
from clu.sources import (
    clu,
    proc_cpuinfo,
    ipmitool,
    lscpu,
    no_salt,
    os_release,
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
        clu.Clu,
        proc_cpuinfo.ProcCpuinfo,
        ipmitool.Ipmitool,
        lscpu.Lscpu,
        no_salt.NoSalt,
        os_release.OsRelease,
        proc_uptime.ProcUptime,
        selinux.Selinux,
        sys_dmi.SysDmi,
        udevadm_ram.UdevadmRam,
        uname.Uname,
        virt_what.VirtWhat,
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

    # def parse(self) -> Facts:
    #     """Parse the facts for macOS (Darwin)."""
    #     facts = super().parse()

    #     # Nothing explicitly says Apple, but we know its apple because Darwin is the OS
    #     facts["os.name"] = "Linux"

    #     return facts

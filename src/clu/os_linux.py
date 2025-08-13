"""Doc Incomplete."""

import logging

from clu import Requires
from clu import Facts
from clu import Provides
from clu.debug import panic
from clu.input import text_program, text_file
from clu.conversions import bytes_to_si, seconds_to_text
from clu.os_generic import (
    provides_clu,
    provides_uname,
    requires_uname,
    parse_uname,
    requires_clu,
    parse_clu,
)

log = logging.getLogger(__name__)


def default_facts_os_linux() -> list:
    return [
        "os.name",
        "os.hostname",
        "os.distro.name",
        "os.distro.version",
        "phy.ram",
        "run.uptime",
        "clu.version",
    ]


def provides_os_linux() -> Provides:
    """Define the provider map for Linux."""
    provides = Provides()

    provides_uname(provides)
    provides_virt_what(provides)
    provides_os_release(provides)
    provides_sys_dmi(provides)
    provides_cpuinfo_flags(provides)
    provides_udevadm_ram(provides)
    provides_lscpu(provides)
    provides_selinux(provides)
    provides_no_salt(provides)
    provides_proc_uptime(provides)
    provides_clu(provides)
    provides_ipmitool(provides)

    return provides


def requires_os_linux() -> Requires:
    """Define the requirements for Linux."""
    requires = Requires()

    requires_uname(requires)
    requires_virt_what(requires)
    requires_os_release(requires)
    requires_sys_dmi(requires)
    requires_cpuinfo_flags(requires)
    requires_udevadm_ram(requires)
    requires_lscpu(requires)
    requires_selinux(requires)
    requires_no_salt(requires)
    requires_proc_uptime(requires)
    requires_clu(requires)
    requires_ipmitool(requires)

    return requires


def parse_os_linux() -> Facts:
    """Parse the facts for Linux."""
    facts = Facts()

    facts["os.name"] = "Linux"

    parse_uname(facts)
    parse_virt_what(facts)
    parse_os_release(facts)
    parse_sys_dmi(facts)
    parse_cpuinfo_flags(facts)
    parse_udevadm_ram(facts)
    parse_lscpu(facts)
    parse_selinux(facts)
    parse_no_salt(facts)
    parse_proc_uptime(facts)
    parse_clu(facts)
    parse_ipmitool(facts)

    return facts






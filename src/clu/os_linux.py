"""Doc Incomplete."""

import logging
import re

from clu.requires import Requires
from clu.facts import Facts
from clu import panic
from clu.readers import read_program, read_file
from clu.conversions import bytes_to_si
from clu.os_generic import (
    requires_uname,
    parse_uname,
    requires_clu,
    parse_clu,
)

log = logging.getLogger(__name__)


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

    return requires


def parse_os_linux() -> Facts:
    """Parse the facts for Linux."""
    facts = Facts()

    facts["os.name"] = "Linux"

    parse_uname(facts)
    parse_virt_what(facts)
    parse_uname(facts)
    parse_sys_dmi(facts)
    parse_cpuinfo_flags(facts)
    parse_udevadm_ram(facts)
    parse_lscpu(facts)
    parse_selinux(facts)
    parse_no_salt(facts)
    parse_proc_uptime(facts)
    parse_clu(facts)

    return facts


def requires_os_release(requires: Requires) -> None:
    requires.files.append("/etc/os-release")


def parse_os_release(facts: Facts) -> None:
    data = read_file("/etc/os-release")
    log.debug(f"{data=}")
    if not data:
        facts["os.distro.name"] = "Unknown/Error"
        facts["os.distro.version"] = "Unknown/Error"
        return
    for line in data.splitlines():
        if "=" not in line:
            continue
        key, value = line.split("=", 1)
        key = key.strip().strip('"')
        value = value.strip().strip('"')
        log.debug(f"{key=}")
        log.debug(f"{value=}")

        if key == "ID":
            facts["os.distro.name"] = value
        elif key == "VERSION_ID":
            facts["os.distro.version"] = value


def requires_sys_dmi(requires: Requires) -> None:
    requires.files.extend([
        "/sys/devices/virtual/dmi/id/sys_vendor",
        "/sys/devices/virtual/dmi/id/product_family",
        "/sys/devices/virtual/dmi/id/product_name",
        "/sys/devices/virtual/dmi/id/product_serial",
        "/sys/devices/virtual/dmi/id/product_uuid",
        "/sys/devices/virtual/dmi/id/chassis_vendor",
        "/sys/devices/virtual/dmi/id/chassis_asset_tag",
    ])


def parse_sys_dmi(facts: Facts) -> None:

    if facts["phy.platform"] != "physical":
        log.info("Not a physical platform, skipping sys.dmi parsing")
        return

    keys = [
        "sys.vendor",
        "sys.model.family",
        "sys.model.name",
        "sys.serial_no",
        "sys.uuid",
        "sys.oem",
        "sys.asset_no",
    ]
    entries = [
        "sys_vendor",
        "product_family",
        "product_name",
        "product_serial",
        "product_uuid",
        "chassis_vendor",
        "chassis_asset_tag",
    ]
    if len(keys) != len(entries):
        log.debug(f"{keys=}")
        log.debug(f"{entries=}")
        panic("parse_sys_dmi: keys and entries length don't match: You can't count")
    for idx in range(len(keys)):
        key = keys[idx]
        entry = entries[idx]
        data = read_file(f"/sys/devices/virtual/dmi/id/{entry}")
        log.debug(f"{data=}")
        log.debug(f"{key=}")
        facts[key] = data.strip()


def requires_udevadm_ram(requires: Requires) -> None:
    requires.programs.append("udevadm info --path /devices/virtual/dmi/id")


def parse_udevadm_ram(facts: Facts) -> None:
    data, rc = read_program("udevadm info --path /devices/virtual/dmi/id")
    if data is None or rc != 0:
        return

    # Find all MEMORY_DEVICE_x_SIZE=number
    raw_sizes = re.findall(r"MEMORY_DEVICE_\d+_SIZE=(\d+)", data)
    log.debug(f"{raw_sizes=}")

    sizes = [int(size) for size in raw_sizes]
    log.debug(f"{sizes=}")

    total = sum(sizes)
    log.debug(f"{total=}")

    bytes_str = bytes_to_si(total)
    log.debug(f"{bytes_str=}")

    facts["phy.ram"] = bytes_str


def requires_virt_what(requires: Requires) -> None:
    requires.programs.append("virt-what")


def parse_virt_what(facts: Facts) -> None:
    if "phy.platform" in facts:
        # already set --skip it
        return

    data, rc = read_program("virt-what")
    if data is None or rc != 0:
        facts["phy.platform"] = "Unknown/Error"
        return
    data = data.strip()
    log.debug(f"{data=}")
    if not data:
        data = "physical"
    facts["phy.platform"] = data


def requires_lscpu(requires: Requires) -> None:
    requires.programs.append("lscpu")

# TODO: clean this up, it is a mess because it didn't translate from the original code well
def parse_lscpu(facts: Facts) -> None:
    regexes = {
        r"^ *Model name: *(.+)": "model",
        r"^ *Vendor ID: *(.+)": "vendor",
        r"^ *Core\(s\) per socket: *(\d+)": "cores_per_socket",
        r"^ *Thread\(s\) per core: *(\d+)": "threads_per_core",
        r"^ *Socket\(s\): *(\d+)": "sockets",
        r"^ *CPU\(s\): *(\d+)": "cpus",
    }
    fields = {}
    attr_keys = ["model", "vendor", "cores", "threads", "sockets"]
    data, rc = read_program("lscpu")
    if data is None or rc != 0:
        return

    for regex, field in regexes.items():
        match = re.search(regex, data, re.MULTILINE)
        value = match.group(1).strip() if match else None
        log.debug(f"{value=}")
        if value is not None:
            fields[field] = value
    log.debug(f"{fields=}")
    try:
        fields["cores"] = str(int(fields["cores_per_socket"]) * int(fields["sockets"]))
        fields["threads"] = str(int(fields["threads_per_core"]) * int(fields["cores"]))
    except Exception:
        pass
    log.debug(f"{fields=}")
    for key in attr_keys:
        value = fields.get(key)
        log.debug(f"{key=}")
        log.debug(f"{value=}")
        facts[f"phy.cpu.{key}"] = value


def __has_flags(check_flags, all_flags):
    check_flags = check_flags.split()
    all_flags = all_flags.split()
    log.debug(f"count is {len(check_flags)} {len(all_flags)}")
    log.debug(f"{check_flags=}")
    log.debug(f"{all_flags=}")
    for flag in check_flags:
        if flag not in all_flags:
            return False
    return True


def requires_cpuinfo_flags(requires: Requires) -> None:
    requires.files.append("/proc/cpuinfo")


def parse_cpuinfo_flags(facts: Facts) -> None:
    # we can always assume uname has been parsed
    if facts["phy.arch"] not in ("x86_64", "amd64"):
        log.info("Not an x86_64/amd64 architecture, skipping cpuinfo flags parsing")
        return

    vers = [
        "lm cmov cx8 fpu fxsr mmx syscall sse2",
        "cx16 lahf_lm popcnt sse4_1 sse4_2 ssse3",
        "avx avx2 bmi1 bmi2 f16c fma abm movbe xsave",
        "avx512f avx512bw avx512cd avx512dq avx512vl",
    ]
    data = read_file("/proc/cpuinfo")
    log.debug(f"data={data.splitlines() if data else []}")
    cpu_flags = ""
    if data:
        for line in data.splitlines():
            if line.startswith("flags"):
                cpu_flags = line.split(":", 1)[1].strip()
                break
    log.debug(f"cpu_flags={cpu_flags.split()}")
    cpu_version = 0
    for idx, v in enumerate(vers):
        if __has_flags(v, cpu_flags):
            cpu_version += 1
        else:
            break
    log.debug(f"idx={idx if data else 0}")
    facts["phy.cpu.arch_version"] = f"x86_64_v{cpu_version}"


def requires_selinux(requires: Requires) -> None:
    requires.programs.extend(["selinuxenabled", "getenforce"])


def parse_selinux(facts: Facts) -> None:
    _, rc = read_program("selinuxenabled")
    # man page: "status 0 if SELinux is enabled and 1 if it is not enabled."
    log.debug(f'rc is {rc}')
    if rc == 0:
        facts["os.selinux.enable"] = "True"
    elif rc == 1:
        facts["os.selinux.enable"] = "False"
    else:
        facts["os.selinux.enable"] = "Unknown/Error"

    data, rc = read_program("getenforce")
    log.debug(f"{data=}")
    facts["os.selinux.mode"] = data.strip() if data else "Unknown/Error"


def requires_no_salt(requires: Requires) -> None:
    requires.files.append("/no_salt")


def parse_no_salt(facts: Facts) -> None:
    data = read_file("/no_salt")
    if data is None:
        facts["salt.no_salt.exists"] = "False"
        return
    else:
        log.debug(f"{data=}")
        facts["salt.no_salt.exists"] = "True"
    if not data.strip():
        data = "UNKNOWN"
    facts["salt.no_salt.reason"] = data.strip()


def requires_proc_uptime(requires: Requires) -> None:
    requires.files.append("/proc/uptime")


def parse_proc_uptime(facts: Facts) -> None:
    data = read_file("/proc/uptime")
    log.debug(f"{data=}")
    if not data:
        return
    uptime_secs = data.split()[0]
    log.debug(f"{uptime_secs=}")
    facts["run.uptime"] = uptime_secs


def requires_ip_addr(requires: Requires) -> None:
    requires.programs.append("ip")


def parse_ip_addr(facts: Facts) -> None:
    data, rc = read_program("ip addr")
    if data is None or rc != 0:
        return
    log.debug(f"{data=}")

# TODO:
# 1. Create parse_ram() function to parse ls_mem and if that fails fallback to udevadm_ram()

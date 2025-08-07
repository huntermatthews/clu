"""Doc Incomplete."""

import logging
import re

from clu.requires import Requires
from clu.facts import Facts
from clu.provides import Provides
from clu import panic
from clu.readers import read_program, read_file
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


def provides_os_release(provides: Provides) -> None:
    provides["os.distro.name"] = parse_os_release
    provides["os.distro.version"] = parse_os_release


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


def provides_sys_dmi(provides: Provides) -> None:
    provides["sys.vendor"] = parse_sys_dmi
    provides["sys.model.family"] = parse_sys_dmi
    provides["sys.model.name"] = parse_sys_dmi
    provides["sys.serial_no"] = parse_sys_dmi
    provides["sys.uuid"] = parse_sys_dmi
    provides["sys.oem"] = parse_sys_dmi
    provides["sys.asset_no"] = parse_sys_dmi


def requires_sys_dmi(requires: Requires) -> None:
    requires.files.extend(
        [
            "/sys/devices/virtual/dmi/id/sys_vendor",
            "/sys/devices/virtual/dmi/id/product_family",
            "/sys/devices/virtual/dmi/id/product_name",
            "/sys/devices/virtual/dmi/id/product_serial",
            "/sys/devices/virtual/dmi/id/product_uuid",
            "/sys/devices/virtual/dmi/id/chassis_vendor",
            "/sys/devices/virtual/dmi/id/chassis_asset_tag",
        ]
    )


def parse_sys_dmi(facts: Facts) -> None:
    if "phy.platform" not in facts:
        parse_virt_what(facts)

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


def provides_udevadm_ram(provides: Provides) -> None:
    provides["phy.ram"] = parse_udevadm_ram


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


def provides_virt_what(provides: Provides) -> None:
    provides["phy.platform"] = parse_virt_what


def requires_virt_what(requires: Requires) -> None:
    requires.programs.append("virt-what")


def parse_virt_what(facts: Facts) -> None:
    if "phy.platform" in facts:
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


def provides_lscpu(provides: Provides) -> None:
    provides["phy.cpu.model"] = parse_lscpu
    provides["phy.cpu.vendor"] = parse_lscpu
    provides["phy.cpu.cores"] = parse_lscpu
    provides["phy.cpu.threads"] = parse_lscpu
    provides["phy.cpu.sockets"] = parse_lscpu


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


def provides_cpuinfo_flags(provides: Provides) -> None:
    provides["phy.cpu.arch_version"] = parse_cpuinfo_flags


def requires_cpuinfo_flags(requires: Requires) -> None:
    requires.files.append("/proc/cpuinfo")


def parse_cpuinfo_flags(facts: Facts) -> None:
    if "phy.arch" not in facts:
        parse_uname(facts)

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


def provides_selinux(provides: Provides) -> None:
    provides["os.selinux.enable"] = parse_selinux
    provides["os.selinux.mode"] = parse_selinux


def requires_selinux(requires: Requires) -> None:
    requires.programs.extend(["selinuxenabled", "getenforce"])


def parse_selinux(facts: Facts) -> None:
    _, rc = read_program("selinuxenabled")
    # man page: "status 0 if SELinux is enabled and 1 if it is not enabled."
    log.debug(f"rc is {rc}")
    if rc == 0:
        facts["os.selinux.enable"] = "True"
    elif rc == 1:
        facts["os.selinux.enable"] = "False"
    else:
        facts["os.selinux.enable"] = "Unknown/Error"

    data, rc = read_program("getenforce")
    log.debug(f"{data=}")
    facts["os.selinux.mode"] = data.strip() if data else "Unknown/Error"


def provides_no_salt(provides: Provides) -> None:
    provides["os.no_salt.exists"] = parse_no_salt
    provides["os.no_salt.reason"] = parse_no_salt


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


def provides_proc_uptime(provides: Provides) -> None:
    provides["run.uptime"] = parse_proc_uptime


def requires_proc_uptime(requires: Requires) -> None:
    requires.files.append("/proc/uptime")


def parse_proc_uptime(facts: Facts) -> None:
    data = read_file("/proc/uptime")
    log.debug(f"{data=}")
    if not data:
        return
    uptime_secs = int(float(data.split()[0]))
    log.debug(f"{uptime_secs=}")
    facts["run.uptime"] = seconds_to_text(uptime_secs)


def provides_ip_addr(provides: Provides) -> None:
    pass


def requires_ip_addr(requires: Requires) -> None:
    requires.programs.append("ip")


def parse_ip_addr(facts: Facts) -> None:
    data, rc = read_program("ip addr")
    if data is None or rc != 0:
        return
    log.debug(f"{data=}")


def provides_ipmitool(provides: Provides) -> None:
    for key in [
        "bmc.ipv4_source",
        "bmc.ipv4_address",
        "bmc.ipv4_mask",
        "bmc.mac_address",
        "bmc.firmware_version",
        "bmc.manufacturer_id",
        "bmc.manufacturer_name",
    ]:
        provides[key] = parse_ipmitool


def requires_ipmitool(requires: Requires) -> None:
    requires.programs.append("ipmitool")


def _parse_ipmitool_lan_print(facts: Facts) -> None:
    regexes = {
        r"^IP Address Source *: (.+)": "bmc.ipv4_source",
        r"^IP Address *: (.+)": "bmc.ipv4_address",
        r"^Subnet Mask *: (.+)": "bmc.ipv4_mask",
        r"^MAC Address *: (.+)": "bmc.mac_address",
    }
    fields = {}
    data, rc = read_program("ipmitool lan print")
    if data is None or rc != 0:
        return

    for regex, field in regexes.items():
        match = re.search(regex, data, re.MULTILINE)
        value = match.group(1).strip() if match else None
        log.debug(f"{value=}")
        if value is not None:
            fields[field] = value
    log.debug(f"{fields=}")

    facts.update(fields)


def _parse_ipmitool_mc_info(facts: Facts) -> None:
    regexes = {
        r"^Firmware Revision *: (.+)": "bmc.firmware_version",
        r"^Manufacturer ID *: (.+)": "bmc.manufacturer_id",
        r"^Manufacturer Name *: (.+)": "bmc.manufacturer_name",
    }
    fields = {}
    data, rc = read_program("ipmitool mc info")
    if data is None or rc != 0:
        return

    for regex, field in regexes.items():
        match = re.search(regex, data, re.MULTILINE)
        value = match.group(1).strip() if match else None
        log.debug(f"{value=}")
        if value is not None:
            fields[field] = value
    log.debug(f"{fields=}")

    facts.update(fields)


def parse_ipmitool(facts: Facts) -> None:
    if "phy.platform" not in facts:
        parse_virt_what(facts)

    if facts["phy.platform"] != "physical":
        log.info("Not a physical platform, skipping bmc. parsing")
        return

    _parse_ipmitool_mc_info(facts)
    _parse_ipmitool_lan_print(facts)

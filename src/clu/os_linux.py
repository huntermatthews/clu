"""Doc Incomplete."""

import re

from clu.requires import add_requires
from clu.facts import add_fact, get_fact, get_all_facts
from clu.debug import trace, debug, debug_var, trace_var, panic
from clu.readers import read_program, read_file
from clu.conversions import bytes_to_si
from clu.os_generic import (
    requires_uname,
#    parse_uname,
    requires_uptime,
    parse_uptime,
    requires_clu,
    parse_clu,
)


def requires_os_linux():
    trace("requires_os_linux begin")

    requires_uname()
    requires_virt_what()
    requires_os_release()
    requires_sys_dmi()
    requires_cpuinfo_flags()
    requires_udevadm_ram()
    requires_lscpu()
    requires_selinux()
    requires_no_salt()
    requires_uptime()
    requires_clu()


def parse_os_linux():
    trace("parse_os_linux begin")
    add_fact("os.name", "Linux")

    # parse_uname() done already
    parse_virt_what()
    parse_os_release()
    parse_sys_dmi()
    parse_cpuinfo_flags()
    parse_udevadm_ram()
    parse_lscpu()
    parse_selinux()
    parse_no_salt()
    parse_uptime()
    parse_clu()


def requires_os_release():
    trace("requires_os_release begin")
    add_requires("files", "/etc/os-release")


def parse_os_release():
    trace("parse_os_release begin")
    data = read_file("/etc/os-release")
    debug_var("data", data)
    if not data:
        return
    for line in data.splitlines():
        if "=" not in line:
            continue
        key, value = line.split("=", 1)
        key = key.strip().strip('"')
        value = value.strip().strip('"')
        debug_var("key", key)
        debug_var("value", value)
        if key == "ID":
            add_fact("os.distro.name", value)
        elif key == "VERSION_ID":
            add_fact("os.distro.version", value)


def requires_sys_dmi():
    trace("requires_sys_dmi begin")
    add_requires("files", "/sys/devices/virtual/dmi/id/sys_vendor")
    add_requires("files", "/sys/devices/virtual/dmi/id/product_family")
    add_requires("files", "/sys/devices/virtual/dmi/id/product_name")
    add_requires("files", "/sys/devices/virtual/dmi/id/product_serial")
    add_requires("files", "/sys/devices/virtual/dmi/id/product_uuid")
    add_requires("files", "/sys/devices/virtual/dmi/id/chassis_vendor")
    add_requires("files", "/sys/devices/virtual/dmi/id/chassis_asset_tag")


def parse_sys_dmi():
    trace("parse_sys_dmi begin")

    parse_virt_what()

    if get_fact("phy.platform") != "physical":
        debug("Not a physical platform, skipping sys.dmi parsing")
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
        debug_var("keys", keys)
        debug_var("entries", entries)
        panic("parse_sys_dmi: keys and entries length don't match: You can't count")
    for idx in range(len(keys)):
        key = keys[idx]
        entry = entries[idx]
        data = read_file(f"/sys/devices/virtual/dmi/id/{entry}")
        debug_var("data", data)
        debug_var("key", key)
        add_fact(key, data.strip())


def requires_udevadm_ram():
    add_requires("programs", "udevadm info --path /devices/virtual/dmi/id")


def parse_udevadm_ram():
    trace("parse_udevadm_ram begin")
    data, rc = read_program("udevadm info --path /devices/virtual/dmi/id")
    if data is None or rc != 0:
        return

    # Find all MEMORY_DEVICE_x_SIZE=number
    raw_sizes = re.findall(r"MEMORY_DEVICE_\d+_SIZE=(\d+)", data)
    debug_var("raw_sizes", raw_sizes)

    sizes = [int(size) for size in raw_sizes]
    debug_var("sizes", sizes)

    total = sum(sizes)
    debug_var("total", total)

    bytes_str = bytes_to_si(total)
    debug_var("bytes", bytes_str)

    add_fact("phy.ram", bytes_str)


def requires_virt_what():
    add_requires("programs", "virt-what")


def parse_virt_what():
    trace("parse_virt_what begin")
    if "phy.platform" in get_all_facts():
        # already set --skip it
        return

    data, rc = read_program("virt-what")
    if data is None or rc != 0:
        add_fact("phy.platform", "Unknown/Error")
        return
    data = data.strip()
    debug_var("data", [data])
    if not data:
        data = "physical"
    add_fact("phy.platform", data)

def requires_lscpu():
    add_requires("programs", "lscpu")


def parse_lscpu():
    trace("parse_lscpu begin")
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
        debug_var("value", value)
        if value is not None:
            fields[field] = value
    debug_var("fields", fields)
    try:
        fields["cores"] = str(int(fields["cores_per_socket"]) * int(fields["sockets"]))
        fields["threads"] = str(int(fields["threads_per_core"]) * int(fields["cores"]))
    except Exception:
        pass
    debug_var("fields", fields)
    for key in attr_keys:
        value = fields.get(key)
        debug_var("key", key)
        debug_var("value", value)
        add_fact(f"phy.cpu.{key}", value)


def __has_flags(check_flags, all_flags):
    trace("__has_flags begin")
    check_flags = check_flags.split()
    all_flags = all_flags.split()
    debug("count is", len(check_flags), len(all_flags))
    debug_var("check_flags", check_flags)
    debug_var("all_flags", all_flags)
    for flag in check_flags:
        if flag not in all_flags:
            return False
    trace("__has_flags end")
    return True


def requires_cpuinfo_flags():
    add_requires("files", "/proc/cpuinfo")


def parse_cpuinfo_flags():
    trace("parse_cpuinfo_flags begin")

    # we can always assume uname has been parsed
    if get_fact("phy.arch") not in ("x86_64", "amd64"):
        debug("Not an x86_64/amd64 architecture, skipping cpuinfo flags parsing")
        return

    vers = [
        "lm cmov cx8 fpu fxsr mmx syscall sse2",
        "cx16 lahf_lm popcnt sse4_1 sse4_2 ssse3",
        "avx avx2 bmi1 bmi2 f16c fma abm movbe xsave",
        "avx512f avx512bw avx512cd avx512dq avx512vl",
    ]
    data = read_file("/proc/cpuinfo")
    trace_var("data", data.splitlines() if data else [])
    cpu_flags = ""
    if data:
        for line in data.splitlines():
            if line.startswith("flags"):
                cpu_flags = line.split(":", 1)[1].strip()
                break
    debug_var("cpu_flags", cpu_flags.split())
    cpu_version = 0
    for idx, v in enumerate(vers):
        if __has_flags(v, cpu_flags):
            cpu_version += 1
        else:
            break
    debug_var("idx", idx if data else 0)
    add_fact("phy.cpu.arch_version", f"x86_64_v{cpu_version}")


def requires_selinux():
    add_requires("programs", "selinuxenabled")
    add_requires("programs", "getenforce")


def parse_selinux():
    trace("parse_selinux begin")
    _, rc = read_program("selinuxenabled")
    # man page: "status 0 if SELinux is enabled and 1 if it is not enabled."
    debug(f'rc is {rc}')
    if rc == 0:
        add_fact("os.selinux.enable", True)
    elif rc == 1:
        add_fact("os.selinux.enable", False)
    else:
        add_fact("os.selinux.enable", "Unknown/Error")

    data, rc = read_program("getenforce")
    debug_var("data", data)
    add_fact("os.selinux.mode", data.strip() if data else "Unknown/Error")


def requires_no_salt():
    add_requires("files", "/no_salt")


def parse_no_salt():
    trace("parse_no_salt begin")
    data = read_file("/no_salt")
    if data is None:
        add_fact("salt.no_salt.exists", False)
        return
    else:
        debug_var("data", data)
        add_fact("salt.no_salt.exists", True)
    if not data.strip():
        data = "UNKNOWN"
    add_fact("salt.no_salt.reason", data.strip())


def requires_proc_uptime():
    add_requires("files", "/proc/uptime")


def parse_proc_uptime():
    trace("parse_proc_uptime begin")
    data = read_file("/proc/uptime")
    debug_var("data", data)
    if not data:
        return
    uptime_secs = data.split()[0]
    debug_var("uptime_secs", uptime_secs)
    add_fact("run.uptime", uptime_secs)


def requires_ip_addr():
    add_requires("programs", "ip")


def parse_ip_addr():
    trace("parse_ip_addr begin")
    data, rc = read_program("ip addr")
    if data is None or rc != 0:
        return
    debug_var("data", data)

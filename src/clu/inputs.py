import re
import os

# Assume ATTRS is a global dictionary
ATTRS = {}


def input_os_release():
    trace("input_os_release begin")
    data = read_file("/etc/os-release")
    debug_var_list("data", data)
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
            ATTRS["os.distro.name"] = value
        elif key == "VERSION_ID":
            ATTRS["os.distro.version"] = value


def input_uname():
    trace("input_uname begin")
    keys = [
        "os.kernel.name",
        "os.hostname",
        "os.kernel.version",
        "phy.arch.name",
        "phy.arch.family",
    ]
    data = read_program("uname", "-snrmp")
    if data is None:
        panic("input_uname: uname command failed")
        return
    data = data.strip().split()
    debug_var_list("data", data)
    if len(keys) != len(data):
        debug("count keys", len(keys))
        debug("count data", len(data))
        debug_var("keys", keys)
        debug_var("data", data)
        panic("input_uname: keys and data length don't match: You can't count")
    for idx in range(len(data)):
        debug_var(f"keys[{idx}]", keys[idx])
        debug_var(f"data[{idx}]", data[idx])
        ATTRS[keys[idx]] = data[idx]


def input_sys_dmi():
    trace("input_sys_dmi begin")
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
        debug("count keys", len(keys))
        debug("count entries", len(entries))
        debug_var("keys", keys)
        debug_var("entries", entries)
        panic("input_sys_dmi: keys and entries length don't match: You can't count")
    for idx in range(len(keys)):
        key = keys[idx]
        entry = entries[idx]
        data = read_file(f"/sys/devices/virtual/dmi/id/{entry}")
        debug_var_list("data", data)
        debug_var("key", key)
        ATTRS[key] = data


def requires_udevadm_ram():
    return "prog:udevadm"


def input_udevadm_ram():
    trace("input_udevadm_ram begin")
    data = read_program("udevadm", "info", "-e")
    if data is None:
        return
    # Find all MEMORY_DEVICE_x_SIZE=number
    raw_sizes = re.findall(r"MEMORY_DEVICE_\d+_SIZE=(\d+)", data)
    debug_var_list("raw_sizes", raw_sizes)
    sizes = [int(size) for size in raw_sizes]
    debug_var_list("sizes", sizes)
    total = sum(sizes)
    debug_var("total", total)
    bytes_str = bytes_to_si(total)
    debug_var("bytes", bytes_str)
    ATTRS["phy.ram.size"] = bytes_str


def requires_virt_what():
    return "prog:virt-what"


def input_virt_what():
    trace("input_virt_what begin")
    data = read_program("virt-what")
    if data is None:
        ATTRS["phy.platform"] = "UNKNOWN"
        return
    data = data.strip()
    debug_var_list("data", [data])
    if not data:
        data = "physical"
    ATTRS["phy.platform"] = data


def requires_lscpu():
    return "prog:lscpu"


def input_lscpu():
    trace("input_lscpu begin")
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
    data = read_program("lscpu")
    if data is None:
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
        ATTRS[f"phy.cpu.{key}"] = value


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
    return "file:/proc/cpuinfo"


def input_cpuinfo_flags():
    trace("input_cpuinfo_flags begin")
    vers = [
        "lm cmov cx8 fpu fxsr mmx syscall sse2",
        "cx16 lahf_lm popcnt sse4_1 sse4_2 ssse3",
        "avx avx2 bmi1 bmi2 f16c fma abm movbe xsave",
        "avx512f avx512bw avx512cd avx512dq avx512vl",
    ]
    data = read_file("/proc/cpuinfo")
    trace_var_list("data", data.splitlines() if data else [])
    cpu_flags = ""
    if data:
        for line in data.splitlines():
            if line.startswith("flags"):
                cpu_flags = line.split(":", 1)[1].strip()
                break
    debug_var_list("cpu_flags", cpu_flags.split())
    cpu_version = 0
    for idx, v in enumerate(vers):
        if __has_flags(v, cpu_flags):
            cpu_version += 1
        else:
            break
    debug_var("idx", idx if data else 0)
    ATTRS["phy.cpu.arch_version"] = f"x86_64_v{cpu_version}"


def requires_selinux():
    return "prog:selinuxenabled prog:getenforce"


def input_selinux():
    trace("input_selinux begin")
    rc = read_program("selinuxenabled")
    if rc is not None:
        ATTRS["os.selinux.enable"] = True
    else:
        ATTRS["os.selinux.enable"] = False
    data = read_program("getenforce")
    debug_var_list("data", [data])
    ATTRS["os.selinux.mode"] = data.strip() if data else None


def requires_gru():
    return ""


def input_gru():
    trace("input_gru begin")
    # These would need to be set appropriately in your Python environment
    ATTRS["gru.binary"] = os.path.realpath(__file__)
    ATTRS["gru.version"] = os.environ.get("_VERSION", "unknown")
    ATTRS["gru.version_info"] = ATTRS["gru.version"].replace(".", " ")
    ATTRS["gru.fish.binary"] = os.environ.get("SHELL", "")
    ATTRS["gru.fish.version"] = os.environ.get("version", "")
    ATTRS["gru.debug_mode"] = os.environ.get("_debug", "")
    ATTRS["gru.path"] = os.environ.get("PATH", "")
    ATTRS["gru.cmdline"] = " ".join(os.sys.argv)


def requires_sw_vers():
    return "prog:sw_vers"


def input_sw_vers():
    trace("input_sw_vers begin")
    data = read_program("sw_vers")
    debug_var_list("data", data.splitlines() if data else [])
    if not data:
        return
    for line in data.splitlines():
        if ":" not in line:
            continue
        key, value = line.split(":", 1)
        key = key.strip()
        value = value.strip()
        debug_var("key", key)
        debug_var("value", value)
        if key == "ProductName":
            ATTRS["os.name"] = value
        elif key == "ProductVersion":
            ATTRS["os.version"] = value
        elif key == "BuildVersion":
            ATTRS["os.build"] = value


def requires_macos_name():
    return ""


def input_macos_name():
    trace("input_macos_name begin")
    version = ATTRS.get("os.version", "")
    major_ver = version.split(".")[0] if version else ""
    debug_var("major_ver", major_ver)
    code_name = None
    if major_ver == "15":
        code_name = "Sequoia"
    elif major_ver == "14":
        code_name = "Sonoma"
    elif major_ver == "13":
        code_name = "Ventura"
    elif major_ver == "12":
        code_name = "Monterey"
    elif major_ver == "11":
        code_name = "Big Sur"
    ATTRS["os.code_name"] = code_name


def requires_no_salt():
    return "file:/no_salt"


def input_no_salt():
    trace("input_no_salt begin")
    data = read_file("/no_salt")
    if data is None:
        ATTRS["salt.no_salt.exists"] = False
        return
    else:
        debug_var_list("data", [data])
        ATTRS["salt.no_salt.exists"] = True
    if not data.strip():
        data = "UNKNOWN"
    ATTRS["salt.no_salt.reason"] = data.strip()


def requires_uptime():
    return "prog:uptime"


def input_uptime():
    trace("input_uptime begin")
    data = read_program("uptime")
    debug_var("data", data)
    if not data:
        return
    match = re.match(r".* up (.*) \d+ user.*", data)
    uptime = match.group(1).rstrip(",") if match else ""
    debug_var("uptime", uptime)
    ATTRS["run.uptime"] = uptime


def requires_proc_uptime():
    return "file:/proc/uptime"


def input_proc_uptime():
    trace("input_proc_uptime begin")
    data = read_file("/proc/uptime")
    debug_var("data", data)
    if not data:
        return
    uptime_secs = data.split()[0]
    debug_var("uptime_secs", uptime_secs)
    ATTRS["run.uptime"] = uptime_secs


def requires_ip_addr():
    return "prog:ip prog:jq"


def input_ip_addr():
    trace("input_ip_addr begin")
    data = read_program("ip", "addr")
    debug_var_list("data", data.splitlines() if data else [])


# Helper stubs for required functions
def trace(msg, *args):
    print(f"TRACE: {msg}", *args)


def debug(msg, *args):
    print(f"DEBUG: {msg}", *args)


def debug_var(var_name, value):
    print(f"DEBUG: {var_name} = {value}")


def debug_var_list(var_name, value_list):
    print(f"DEBUG: {var_name} = {value_list}")


def panic(msg):
    raise RuntimeError(msg)


def read_file(fname):
    try:
        with open(fname, "r") as f:
            return f.read()
    except Exception:
        return None


def read_program(*args):
    import subprocess

    try:
        result = subprocess.run(args, capture_output=True, text=True, check=True)
        return result.stdout
    except Exception:
        return None


def bytes_to_si(size):
    units = ["B", "KB", "MB", "GB", "TB", "PB", "EB"]
    size = float(size)
    for unit in units:
        if size < 1024:
            return f"{size:.1f} {unit}"
        size /= 1024
    return

"""Doc Incomplete."""

from clu import facts
from clu.debug import trace, debug, debug_var, debug_var_list, panic
from clu.old_readers import read_program
from clu.parse_generic import (
    requires_uname,
    parse_uname,
    requires_uptime,
    parse_uptime,
    requires_clu,
    parse_clu,
)


def requires_os_darwin():
    trace("requires_os_darwin begin")
    requires_uname()
    requires_sw_vers()
    requires_macos_name()
    requires_uptime()
    requires_clu()


def parse_os_darwin():
    trace("os_darwin_parse begin")

    # Nothing explicitly says Apple, but we know its apple because Darwin is the OS
    facts["sys.vendor"] = "Apple"
    parse_uname()
    parse_sw_vers()
    parse_macos_name()
    parse_uptime()
    parse_clu()


def requires_sw_vers():
    return "prog:sw_vers"


def parse_sw_vers():
    trace("parse_sw_vers begin")
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
            facts["os.name"] = value
        elif key == "ProductVersion":
            facts["os.version"] = value
        elif key == "BuildVersion":
            facts["os.build"] = value


def requires_macos_name():
    return ""


def parse_macos_name():
    trace("parse_macos_name begin")
    version = facts.get("os.version", "")
    major_ver = version.split(".")[0] if version else ""
    debug_var("major_ver", major_ver)
    if not major_ver:
        panic("parse_macos_name: os.version is not set or empty")

    if major_ver == "26":
        code_name = "Tahoe"
    # BUG: check for 16-25 and error out
    elif major_ver == "15":
        code_name = "Sequoia"
    elif major_ver == "14":
        code_name = "Sonoma"
    elif major_ver == "13":
        code_name = "Ventura"
    elif major_ver == "12":
        code_name = "Monterey"
    elif major_ver == "11":
        code_name = "Big Sur"
    facts["os.code_name"] = code_name

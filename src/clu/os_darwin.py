"""Doc Incomplete."""

from clu import facts, requires
from clu.debug import trace, debug_var, panic
from clu.readers import read_program
from clu.os_generic import (
    requires_uname,
#     parse_uname,
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

    # parse_uname() done already
    parse_sw_vers()
    parse_macos_name()
    parse_uptime()
    parse_clu()


def requires_sw_vers():
    requires["programs"].append("sw_vers")


def parse_sw_vers():
    trace("parse_sw_vers begin")
    data, rc = read_program("sw_vers")
    debug_var("data", data)
    if data is None or rc != 0:
        panic("parse_sw_vers: uname command failed")
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
    # its just logic code - there are no external requirements for this
    pass


def parse_macos_name():
    trace("parse_macos_name begin")
    version = facts.get("os.version", "")
    if not version:
        panic("parse_macos_name: os.version is not set or empty")

    major_ver = version.split(".")[0]
    debug_var("major_ver", major_ver)

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
    else:
        # Note that for older than 11, the logic of the code name changes
        # and thats WAY out of support for us
        code_name = f"Unknown-{major_ver}"
    facts["os.code_name"] = code_name


## TODO: add more requirements and parsing functions for macOS
## 1. system_profiler SPSoftwareDataType -json

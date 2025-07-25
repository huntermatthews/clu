"""Doc Incomplete."""

from clu.requires import Requires
from clu.facts import Facts
from clu.debug import debug_var, panic
from clu.readers import read_program
from clu.os_generic import (
    requires_uname,
#    parse_uname,
    requires_uptime,
    parse_uptime,
    requires_clu,
    parse_clu,
)


def requires_os_darwin(requires: Requires) -> None:
    requires_uname(requires)
    requires_sw_vers(requires)
    requires_macos_name(requires)
    requires_uptime(requires)
    requires_clu(requires)
    requires_systemversion_plist()


def parse_os_darwin(facts: Facts) -> None:
    # Nothing explicitly says Apple, but we know its apple because Darwin is the OS
    facts["sys.vendor"] = "Apple"

    # parse_uname() done already
    parse_sw_vers(facts)
    parse_macos_name(facts)
    parse_uptime(facts)
    parse_clu(facts)


def requires_sw_vers(requires: Requires) -> None:
    requires.programs.append("sw_vers")


def parse_sw_vers(facts: Facts) -> None:
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


def requires_macos_name(requires: Requires) -> None:
    # its just logic code - there are no external requirements for this
    pass


def parse_macos_name(facts: Facts) -> None:
    version = facts["os.version"]
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


def requires_systemversion_plist(requires: Requires) -> None:
    requires.files.append("/System/Library/CoreServices/SystemVersion.plist")


## TODO: add more requirements and parsing functions for macOS
## 1. system_profiler SPSoftwareDataType -json

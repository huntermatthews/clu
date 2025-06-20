"""Doc Incomplete."""

import os
import re
import sys


from clu import facts, config
from clu.debug import trace, debug, debug_var, debug_var_list, panic
from clu.readers import read_program

# from clu.parse_linux import parse_os_release


def requires_os_test():
    trace("requires_os_test begin")
    print("not implemented")


def parse_os_test():
    trace("parse_os_test begin")
    print("parse_os_release() here I am")


def requires_os_unsupported():
    trace("requires_os_unsupported begin")
    panic("Unsupported OS")


def parse_os_unsupported():
    trace("parse_os_unsupported begin")
    panic("Unsupported OS")


def requires_uname():
    trace("requires_uname begin")
    return "prog:uname"


def parse_uname():
    trace("parse_uname begin")
    print(f"{config=}")
    keys = [
        "os.kernel.name",
        "os.hostname",
        "os.kernel.version",
        "phy.arch.name",
        "phy.arch.family",
    ]
    data = read_program("uname", "-snrmp")
    debug_var("data", data)
    debug("data type", type(data))
    if data is None:
        panic("parse_uname: uname command failed")
        return
    data = data.strip().split()
    debug_var_list("data", data)
    if len(keys) != len(data):
        debug("count keys", len(keys))
        debug("count data", len(data))
        debug_var("keys", keys)
        debug_var("data", data)
        panic("parse_uname: keys and data length don't match: You can't count")
    for idx in range(len(data)):
        debug_var(f"keys[{idx}]", keys[idx])
        debug_var(f"data[{idx}]", data[idx])
        facts[keys[idx]] = data[idx]


def requires_clu():
    return ""


def parse_clu():
    trace("parse_clu begin")

    facts["clu.binary"] = os.path.realpath(__file__)
    # facts["clu.version"] = os.environ.get("_VERSION", "unknown")
    # facts["clu.version_info"] = facts["clu.version"].replace(".", " ")
    # facts["clu.python.binary"] = os.environ.get("SHELL", "")
    # facts["clu.python.version"] = os.environ.get("version", "")
    # facts["clu.debug_mode"] = os.environ.get("_debug", "")
    facts["clu.path"] = os.environ.get("PATH", "")
    facts["clu.cmdline"] = " ".join(sys.argv)


def requires_uptime():
    return "prog:uptime"


def parse_uptime():
    trace("parse_uptime begin")
    data = read_program("uptime")
    debug_var("data", data)
    if not data:
        return
    match = re.match(r".* up (.*) \d+ user.*", data)
    uptime = match.group(1).rstrip(",") if match else ""
    debug_var("uptime", uptime)
    facts["run.uptime"] = uptime

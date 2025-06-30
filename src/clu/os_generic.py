"""Doc Incomplete."""

import os
import re
import sys


from clu import facts
from clu.debug import trace, debug_var, trace_var, panic
from clu.readers import read_program


def requires_uname():
    trace("requires_uname begin")
    return "prog:uname"


def parse_uname():
    trace("parse_uname begin")

    keys = [
        "os.kernel.name",
        "os.hostname",
        "os.kernel.version",
        "phy.arch",
    ]
    data, rc = read_program("uname -snrm")
    trace_var("data", data)
    trace_var("rc", rc)
    if data is None or rc != 0:
        panic("parse_uname: uname command failed")
    data = data.strip().split()
    if len(keys) != len(data):
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
    # print('getcwd:      ', os.getcwd())
    # print('__file__:    ', __file__)

def requires_uptime():
    return "prog:uptime"


def parse_uptime():
    trace("parse_uptime begin")
    data, rc = read_program("uptime")
    if data is None or rc != 0:
        panic("parse_uptime: uptime command failed")
    debug_var("data", data)
    match = re.match(r".*up *(.*) \d+ user.*", data)
    uptime = match.group(1).rstrip(",") if match else ""
    debug_var("uptime", uptime)
    facts["run.uptime"] = uptime

"""Doc Incomplete."""

import os
import re
import sys


from clu import __about__
from clu.facts import add_fact
from clu.requires import add_requires
from clu.debug import debug_var, trace_var, panic
from clu.readers import read_program


def requires_uname():
    add_requires("programs", "uname -snrm")


def parse_uname():
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
        add_fact(keys[idx], data[idx])


def requires_clu():
    # No specific requirements for clu group
    pass

def parse_clu():
    add_fact("clu.binary", sys.argv[0])
    add_fact("clu.version", __about__.__version__)
    add_fact("clu.python.binary", sys.executable)
    add_fact("clu.python.version", ".".join(map(str, sys.version_info[:3])))
#    facts["clu.path"] = os.environ.get("PATH", "")
    add_fact("clu.cmdline", " ".join(sys.argv))
    add_fact("clu.cwd", os.getcwd())


def requires_uptime():
    add_requires("programs", "uptime")


def parse_uptime():
    data, rc = read_program("uptime")
    if data is None or rc != 0:
        panic("parse_uptime: uptime command failed")
    debug_var("data", data)
    match = re.match(r".*up *(.*) \d+ user.*", data)
    uptime = match.group(1).rstrip(",") if match else "unknown / error"
    debug_var("uptime", uptime)
    add_fact("run.uptime", uptime)

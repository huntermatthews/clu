"""Doc Incomplete."""

import os
import re
import sys


from clu import __about__
from clu.facts import Facts
from clu.requires import Requires
from clu.debug import debug_var, trace_var, panic
from clu.readers import read_program


def requires_uname(requires: Requires) -> None:
    requires.programs.append("uname -snrm")


def parse_uname(facts: Facts) -> None:
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


def requires_clu(requires: Requires) -> None:
    # No specific requirements for clu group
    pass

def parse_clu(facts: Facts) -> None:
    facts["clu.binary"] = sys.argv[0]
    facts["clu.version"] = __about__.__version__
    facts["clu.python.binary"] = sys.executable
    facts["clu.python.version"] = ".".join(map(str, sys.version_info[:3]))
#    facts["clu.path"] = os.environ.get("PATH", "")
    facts["clu.cmdline"] = " ".join(sys.argv)
    facts["clu.cwd"] = os.getcwd()


def requires_uptime(requires: Requires) -> None:
    requires.programs.append("uptime")


def parse_uptime(facts: Facts) -> None:
    data, rc = read_program("uptime")
    if data is None or rc != 0:
        panic("parse_uptime: uptime command failed")
    debug_var("data", data)
    match = re.match(r".*up *(.*) \d+ user.*", data)
    uptime = match.group(1).rstrip(",") if match else "unknown / error"
    debug_var("uptime", uptime)
    facts["run.uptime"] = uptime

"""Doc Incomplete."""

import logging
import os
import re
import sys


from clu import __about__
from clu.facts import Facts
from clu.requires import Requires
from clu import panic
from clu.readers import read_program

log = logging.getLogger(__name__)


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
    log.debug(f"{data=}")
    log.debug(f"{rc=}")
    if data is None or rc != 0:
        panic("parse_uname: uname command failed")
    data = data.strip().split()
    if len(keys) != len(data):
        log.debug(f"{keys=}")
        log.debug(f"{data=}")
        panic("parse_uname: keys and data length don't match: You can't count")

    for idx in range(len(data)):
        log.debug(f"{keys[idx]=}")
        log.debug(f"{data[idx]=}")
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
    log.debug(f"{data=}")
    match = re.match(r".*up *(.*) \d+ user.*", data)
    uptime = match.group(1).rstrip(",") if match else "unknown / error"
    log.debug(f"{uptime=}")
    facts["run.uptime"] = uptime

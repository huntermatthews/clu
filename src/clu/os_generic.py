"""Doc Incomplete."""

import logging
import os
import re
import sys
import datetime
import getpass

from clu import __about__
from clu import Facts
from clu import Provides
from clu import Requires
from clu.debug import panic
from clu.input import text_program

log = logging.getLogger(__name__)


def provides_uname(provides: Provides) -> None:
    provides["os.kernel.name"] = parse_uname
    provides["os.hostname"] = parse_uname
    provides["os.kernel.version"] = parse_uname
    provides["phy.arch"] = parse_uname


def requires_uname(requires: Requires) -> None:
    requires.programs.append("uname -snrm")


def parse_uname(facts: Facts) -> None:
    if "phy.arch" in facts:
        return

    keys = [
        "os.kernel.name",
        "os.hostname",
        "os.kernel.version",
        "phy.arch",
    ]
    data, rc = text_program("uname -snrm")
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


def provides_clu(provides: Provides) -> None:
    provides["clu.binary"] = parse_clu
    provides["clu.version"] = parse_clu
    provides["clu.python.binary"] = parse_clu
    provides["clu.python.version"] = parse_clu
    provides["clu.cmdline"] = parse_clu
    provides["clu.cwd"] = parse_clu
    provides["clu.user"] = parse_clu
    provides["clu.date"] = parse_clu


def requires_clu(requires: Requires) -> None:
    # No specific requirements for clu group
    pass


def _get_rfc3339_timestamp() -> str:
    """Get the current time in RFC 3339 format."""
    now_utc = datetime.datetime.now(datetime.timezone.utc)
    return now_utc.isoformat(sep="T", timespec="seconds")


def parse_clu(facts: Facts) -> None:
    facts["clu.binary"] = sys.argv[0]
    facts["clu.version"] = __about__.__version__
    facts["clu.python.binary"] = sys.executable
    facts["clu.python.version"] = ".".join(map(str, sys.version_info[:3]))
    facts["clu.cmdline"] = " ".join(sys.argv)
    facts["clu.cwd"] = os.getcwd()
    facts["clu.user"] = getpass.getuser()
    facts["clu.date"] = _get_rfc3339_timestamp()


def provides_uptime(provides: Provides) -> None:
    provides["run.uptime"] = parse_uptime


def requires_uptime(requires: Requires) -> None:
    requires.programs.append("uptime")


def parse_uptime(facts: Facts) -> None:
    data, rc = text_program("uptime")
    if data is None or rc != 0:
        panic("parse_uptime: uptime command failed")
    log.debug(f"{data=}")
    match = re.match(r".*up *(.*) \d+ user.*", data)
    uptime = match.group(1).rstrip(",") if match else "unknown / error"
    log.debug(f"{uptime=}")
    facts["run.uptime"] = uptime

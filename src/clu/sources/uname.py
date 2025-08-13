import logging

from clu import Facts, Provides, Requires, panic
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

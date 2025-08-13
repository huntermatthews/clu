import logging
import re

from clu import Facts, Provides, Requires, panic
from clu.input import text_program

log = logging.getLogger(__name__)


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

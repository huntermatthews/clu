import logging

from clu import Facts, Provides, Requires
from clu.input import text_program

log = logging.getLogger(__name__)


def provides_sw_vers(provides: Provides) -> None:
    provides["os.name"] = parse_sw_vers
    provides["os.version"] = parse_sw_vers
    provides["os.build"] = parse_sw_vers


def requires_sw_vers(requires: Requires) -> None:
    requires.programs.append("sw_vers")


def parse_sw_vers(facts: Facts) -> None:
    if "os.version" in facts:
        # we already ran...
        return
    data, rc = text_program("sw_vers")
    log.debug(f"{data=}")
    if data is None or rc != 0:
        return
    for line in data.splitlines():
        if ":" not in line:
            continue
        key, value = line.split(":", 1)
        key = key.strip()
        value = value.strip()
        log.debug(f"{key=}")
        log.debug(f"{value=}")

        if key == "ProductName":
            facts["os.name"] = value
        elif key == "ProductVersion":
            facts["os.version"] = value
        elif key == "BuildVersion":
            facts["os.build"] = value

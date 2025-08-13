import logging

from clu import Facts, Provides, Requires
from clu.input import text_program

log = logging.getLogger(__name__)


def provides_virt_what(provides: Provides) -> None:
    provides["phy.platform"] = parse_virt_what


def requires_virt_what(requires: Requires) -> None:
    requires.programs.append("virt-what")


def parse_virt_what(facts: Facts) -> None:
    if "phy.platform" in facts:
        return

    data, rc = text_program("virt-what")
    if data is None or rc != 0:
        facts["phy.platform"] = "Unknown/Error"
        return
    data = data.strip()
    log.debug(f"{data=}")
    if not data:
        data = "physical"
    facts["phy.platform"] = data

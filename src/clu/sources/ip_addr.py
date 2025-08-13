import logging

from clu import Facts, Provides, Requires
from clu.input import text_program

log = logging.getLogger(__name__)


def provides_ip_addr(provides: Provides) -> None:
    pass


def requires_ip_addr(requires: Requires) -> None:
    requires.programs.append("ip")


def parse_ip_addr(facts: Facts) -> None:
    data, rc = text_program("ip addr")
    if data is None or rc != 0:
        return
    log.debug(f"{data=}")

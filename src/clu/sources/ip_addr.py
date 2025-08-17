import logging

from clu import Facts, Provides, Requires, Source
from clu.input import text_program

log = logging.getLogger(__name__)


class IpAddr(Source):
    def provides(self, provides: Provides) -> None:
        for key in [
            "ip.addr",
            "ip.mask",
            "ip.gateway",
        ]:
            provides[key] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("ip")

    def parse(self, facts: Facts) -> None:
        data, rc = text_program("ip addr")
        if data is None or rc != 0:
            return
        log.debug(f"{data=}")

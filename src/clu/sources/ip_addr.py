import logging

from clu import Facts, Provides, Requires, Source

from clu.input import text_program

log = logging.getLogger(__name__)


class IpAddr(Source):
    def provides(self) -> Provides:
        provides = Provides()
        for key in [
            "ip.addr",
            "ip.mask",
            "ip.gateway",
        ]:
            provides[key] = self
        return provides

    def requires(self) -> Requires:
        requires = Requires()
        requires.programs.append("ip")
        return requires

    def parse(self, facts: Facts) -> Facts:
        data, rc = text_program("ip addr")
        if data is None or rc != 0:
            return facts
        log.debug(f"{data=}")
        return facts

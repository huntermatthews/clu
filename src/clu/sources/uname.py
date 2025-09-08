import logging

from clu import facts, Tier, Provides, Requires
from clu.input import text_program
from clu.sources import Source, PARSE_FAIL_MSG

log = logging.getLogger(__name__)


class Uname(Source):
    _names = {
        "os.kernel.name": Tier.PRIMARY,
        "os.hostname": Tier.PRIMARY,
        "os.kernel.version": Tier.OTHER,
        "phy.arch": Tier.PRIMARY,
    }

    def provides(self, provides: Provides) -> None:
        for name in self._names:
            provides[name] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("uname -snrm")

    def parse(self) -> None:
        if "os.kernel.name" in facts:
            return

        data, rc = text_program("uname -snrm")
        log.debug(f"{data=}")
        log.debug(f"{rc=}")
        if data == "" or rc != 0:
            for name in self._names:
                facts.add(self._names[name], name, PARSE_FAIL_MSG)
            return

        if data:
            for name, value in zip(self._names, data.strip().split()):
                facts.add(self._names[name], name, value)

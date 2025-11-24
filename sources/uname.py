import logging

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts
from clu.sources import Source
from clu.input import text_program
from clu.sources import PARSE_FAIL_MSG

log = logging.getLogger(__name__)


class Uname(Source):
    _keys = [
        "os.kernel.name",
        "os.hostname",
        "os.kernel.version",
        "phy.arch",
    ]

    def provides(self, provides: Provides) -> None:
        for key in self._keys:
            provides[key] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("uname -snrm")

    def parse(self, facts: Facts) -> None:
        if "os.kernel.name" in facts:
            return

        data, rc = text_program("uname -snrm")
        log.debug(f"{data=}")
        log.debug(f"{rc=}")
        if data == "" or rc != 0:
            for key in self._keys:
                facts[key] = PARSE_FAIL_MSG
            return

        if data:
            for key, value in zip(self._keys, data.strip().split()):
                facts[key] = value

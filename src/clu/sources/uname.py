import logging

from clu import Facts, Provides, Requires, Source
from clu.input import text_program

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
        data, rc = text_program("uname -snrm")
        log.trace(f"{data=}")
        log.trace(f"{rc=}")
        if data == "" or rc != 0:
            for key in self._keys:
                facts[key] = "Unknown/Error"
            return

        if data:
            for key, value in zip(self._keys, data.strip().split()):
                facts[key] = value

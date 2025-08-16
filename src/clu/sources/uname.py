import logging

from clu import Facts, Provides, Requires, Source
from clu.debug import panic
from clu.input import text_program

log = logging.getLogger(__name__)


class Uname(Source):
    def provides(self, provides: Provides) -> None:
        provides["os.kernel.name"] = self
        provides["os.hostname"] = self
        provides["os.kernel.version"] = self
        provides["phy.arch"] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("uname -snrm")

    def parse(self, facts: Facts) -> None:
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
        if data is not None:
            for key, value in zip(keys, data.strip().split()):
                facts[key] = value

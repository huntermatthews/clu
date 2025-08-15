import logging

from clu import Facts, Provides, Requires, Source
from clu.debug import panic
from clu.input import text_program

log = logging.getLogger(__name__)


class Uname(Source):
    def provides(self) -> Provides:
        provides = Provides()
        provides["os.kernel.name"] = self.parse
        provides["os.hostname"] = self.parse
        provides["os.kernel.version"] = self.parse
        provides["phy.arch"] = self.parse
        return provides

    def requires(self) -> Requires:
        requires = Requires()
        requires.programs.append("uname -snrm")
        return requires

    def parse(self) -> Facts:
        facts = Facts()

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

        for key, value in zip(keys, data.strip().split()):
            facts[key] = value

        return facts

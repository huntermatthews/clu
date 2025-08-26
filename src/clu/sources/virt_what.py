import logging

from clu import Facts, Provides, Requires, Source
from clu.input import text_program

log = logging.getLogger(__name__)


class VirtWhat(Source):
    def provides(self, provides: Provides) -> None:
        provides["phy.platform"] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("virt-what")

    def parse(self, facts: Facts) -> None:
        data, rc = text_program("virt-what")
        log.trace(f"{rc=}")
        if rc != 0:
            facts["phy.platform"] = "Unknown/Error"
            return
        data = data.strip()
        log.trace(f"{data=}")
        if not data:
            data = "physical"
        facts["phy.platform"] = data

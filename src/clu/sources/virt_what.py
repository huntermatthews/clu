import logging

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts
from clu.sources import Source
from clu.input import text_program
from clu.sources import PARSE_FAIL_MSG

log = logging.getLogger(__name__)


class VirtWhat(Source):
    def provides(self, provides: Provides) -> None:
        provides["phy.platform"] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("virt-what")

    def parse(self, facts: Facts) -> None:
        if "phy.platform" in facts:
            return

        data, rc = text_program("virt-what")
        log.debug(f"{rc=}")
        if rc != 0:
            facts["phy.platform"] = PARSE_FAIL_MSG
            return
        data = data.strip()
        data = data.replace("\n", ", ")
        log.debug(f"{data=}")
        if not data:
            data = "physical"
        facts["phy.platform"] = data

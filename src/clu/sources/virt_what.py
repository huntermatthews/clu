import logging

from clu import Facts, Provides, Requires, Source
from clu.input import text_program

log = logging.getLogger(__name__)


class VirtWhat(Source):
    def provides(self) -> Provides:
        provides = Provides()
        provides["phy.platform"] = self
        return provides

    def requires(self) -> Requires:
        requires = Requires()
        requires.programs.append("virt-what")
        return requires

    def parse(self, facts: Facts) -> Facts:
        if "phy.platform" in facts:
            return facts

        data, rc = text_program("virt-what")
        if data is None or rc != 0:
            facts["phy.platform"] = "Unknown/Error"
            return facts
        data = data.strip()
        log.debug(f"{data=}")
        if not data:
            data = "physical"
        facts["phy.platform"] = data
        return facts

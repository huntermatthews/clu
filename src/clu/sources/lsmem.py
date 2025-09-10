import logging

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts
from clu.sources import Source
from clu.input import text_program
from clu.sources import PARSE_FAIL_MSG
from clu.conversions import bytes_to_si

log = logging.getLogger(__name__)


class Lsmem(Source):
    def provides(self, provides: Provides) -> None:
        provides["phy.ram"] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("lsmem --summary --bytes")

    def parse(self, facts: Facts) -> None:
        data, rc = text_program("lsmem --summary --bytes")
        log.debug(f"{data=}")
        if not data or rc != 0:
            facts["phy.ram"] = PARSE_FAIL_MSG
            return

        for line in data.splitlines():
            if line.startswith("Total online memory"):
                byte_count = line.split(":", 1)[1]
                byte_count = byte_count.strip()
                log.debug(f"{byte_count=}")

                facts["phy.ram"] = bytes_to_si(float(byte_count))

                # TODO: Once we have verbosity again...
                # facts["phy.ram.used"] = "0"
                # facts["phy.ram.free"] = value

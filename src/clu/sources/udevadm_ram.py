import logging
import re

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts
from clu.sources import Source
from clu.input import text_program
from clu.conversions import bytes_to_si
from clu.sources import PARSE_FAIL_MSG

log = logging.getLogger(__name__)


class UdevadmRam(Source):
    def provides(self, provides: Provides) -> None:
        provides["phy.ram"] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("udevadm info --path /devices/virtual/dmi/id")

    def parse(self, facts: Facts) -> None:
        data, rc = text_program("udevadm info --path /devices/virtual/dmi/id")
        if data == "" or rc != 0:
            facts["phy.ram"] = PARSE_FAIL_MSG
            return

        # Find all MEMORY_DEVICE_x_SIZE=number
        raw_sizes = re.findall(r"MEMORY_DEVICE_\d+_SIZE=(\d+)", data)
        log.debug(f"{raw_sizes=}")

        sizes = [int(size) for size in raw_sizes]
        log.debug(f"{sizes=}")

        total = sum(sizes)
        log.debug(f"{total=}")

        bytes_str = bytes_to_si(total)
        log.debug(f"{bytes_str=}")

        facts["phy.ram"] = bytes_str

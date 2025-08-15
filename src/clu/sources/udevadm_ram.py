import logging
import re

from clu import Facts, Provides, Requires, Source
from clu.input import text_program
from clu.conversions import bytes_to_si

log = logging.getLogger(__name__)


class UdevadmRam(Source):
    def provides(self) -> Provides:
        provides = Provides()
        provides["phy.ram"] = self
        return provides

    def requires(self) -> Requires:
        requires = Requires()
        requires.programs.append("udevadm info --path /devices/virtual/dmi/id")
        return requires

    def parse(self, facts: Facts) -> Facts:
        data, rc = text_program("udevadm info --path /devices/virtual/dmi/id")
        if data is None or rc != 0:
            return facts

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

        return facts

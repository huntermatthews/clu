import logging
import re

from clu import Facts, Provides, Requires
from clu.input import text_program
from clu.conversions import bytes_to_si

log = logging.getLogger(__name__)


def provides_udevadm_ram(provides: Provides) -> None:
    provides["phy.ram"] = parse_udevadm_ram


def requires_udevadm_ram(requires: Requires) -> None:
    requires.programs.append("udevadm info --path /devices/virtual/dmi/id")


def parse_udevadm_ram(facts: Facts) -> None:
    data, rc = text_program("udevadm info --path /devices/virtual/dmi/id")
    if data is None or rc != 0:
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

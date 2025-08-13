import logging
import re

from clu import Facts, Provides, Requires
from clu.input import text_program

log = logging.getLogger(__name__)


def provides_ipmitool(provides: Provides) -> None:
    for key in [
        "bmc.ipv4_source",
        "bmc.ipv4_address",
        "bmc.ipv4_mask",
        "bmc.mac_address",
        "bmc.firmware_version",
        "bmc.manufacturer_id",
        "bmc.manufacturer_name",
    ]:
        provides[key] = parse_ipmitool


def requires_ipmitool(requires: Requires) -> None:
    requires.programs.append("ipmitool")


def _parse_ipmitool_lan_print(facts: Facts) -> None:
    regexes = {
        r"^IP Address Source *: (.+)": "bmc.ipv4_source",
        r"^IP Address *: (.+)": "bmc.ipv4_address",
        r"^Subnet Mask *: (.+)": "bmc.ipv4_mask",
        r"^MAC Address *: (.+)": "bmc.mac_address",
    }
    fields = {}
    data, rc = text_program("ipmitool lan print")
    if data is None or rc != 0:
        return

    for regex, field in regexes.items():
        match = re.search(regex, data, re.MULTILINE)
        value = match.group(1).strip() if match else None
        log.debug(f"{value=}")
        if value is not None:
            fields[field] = value
    log.debug(f"{fields=}")

    facts.update(fields)


def _parse_ipmitool_mc_info(facts: Facts) -> None:
    regexes = {
        r"^Firmware Revision *: (.+)": "bmc.firmware_version",
        r"^Manufacturer ID *: (.+)": "bmc.manufacturer_id",
        r"^Manufacturer Name *: (.+)": "bmc.manufacturer_name",
    }
    fields = {}
    data, rc = text_program("ipmitool mc info")
    if data is None or rc != 0:
        return

    for regex, field in regexes.items():
        match = re.search(regex, data, re.MULTILINE)
        value = match.group(1).strip() if match else None
        log.debug(f"{value=}")
        if value is not None:
            fields[field] = value
    log.debug(f"{fields=}")

    facts.update(fields)


def parse_ipmitool(facts: Facts) -> None:
    if "phy.platform" not in facts:
        parse_virt_what(facts)

    if facts["phy.platform"] != "physical":
        log.info("Not a physical platform, skipping bmc. parsing")
        return

    _parse_ipmitool_mc_info(facts)
    _parse_ipmitool_lan_print(facts)

import logging
import re

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts
from clu.sources import Source, PARSE_FAIL_MSG
from clu.input import text_program

log = logging.getLogger(__name__)


class Ipmitool(Source):
    _lan_print_keys = [
        "bmc.ipv4_source",
        "bmc.ipv4_address",
        "bmc.ipv4_mask",
        "bmc.mac_address",
    ]
    _mc_info_keys = [
        "bmc.firmware_version",
        "bmc.manufacturer_id",
        "bmc.manufacturer_name",
    ]

    def provides(self, provides: Provides) -> None:
        for key in self._lan_print_keys + self._mc_info_keys:
            provides[key] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("ipmitool")

    def parse(self, facts: Facts) -> None:
        if facts["phy.platform"] != "physical":
            log.info("Not a physical platform, skipping bmc parsing")
            return

        self._parse_ipmitool_mc_info(facts)
        self._parse_ipmitool_lan_print(facts)

    def _parse_ipmitool_lan_print(self, facts: Facts) -> None:
        regexes = {
            r"^IP Address Source *: (.+)": "bmc.ipv4_source",
            r"^IP Address *: (.+)": "bmc.ipv4_address",
            r"^Subnet Mask *: (.+)": "bmc.ipv4_mask",
            r"^MAC Address *: (.+)": "bmc.mac_address",
        }
        fields = {}
        data, rc = text_program("ipmitool lan print")
        if data == "" or rc != 0:
            log.warning("Failed to run ipmitool lan print")
            for name in self._lan_print_keys:
                facts[name] = PARSE_FAIL_MSG
            return

        for regex, field in regexes.items():
            match = re.search(regex, data, re.MULTILINE)
            value = match.group(1).strip() if match else None
            log.debug(f"{value=}")
            if value is not None:
                fields[field] = value
        log.debug(f"{fields=}")

        facts.update(fields)

    def _parse_ipmitool_mc_info(self, facts: Facts) -> None:
        regexes = {
            r"^Firmware Revision *: (.+)": "bmc.firmware_version",
            r"^Manufacturer ID *: (.+)": "bmc.manufacturer_id",
            r"^Manufacturer Name *: (.+)": "bmc.manufacturer_name",
        }
        fields = {}
        data, rc = text_program("ipmitool mc info")
        if data == "" or rc != 0:
            log.warning("Failed to run ipmitool mc info")
            for name in self._mc_info_keys:
                facts[name] = PARSE_FAIL_MSG
            return

        for regex, field in regexes.items():
            match = re.search(regex, data, re.MULTILINE)
            value = match.group(1).strip() if match else None
            log.debug(f"{value=}")
            if value is not None:
                fields[field] = value
        log.debug(f"{fields=}")

        facts.update(fields)

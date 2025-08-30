import logging
import re

from clu import Facts, Provides, Requires, Source
from clu.input import text_program

log = logging.getLogger(__name__)


class Ipmitool(Source):
    def provides(self, provides: Provides) -> None:
        for key in [
            "bmc.ipv4_source",
            "bmc.ipv4_address",
            "bmc.ipv4_mask",
            "bmc.mac_address",
            "bmc.firmware_version",
            "bmc.manufacturer_id",
            "bmc.manufacturer_name",
        ]:
            provides[key] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("ipmitool")

    def parse(self, facts: Facts) -> None:
        if facts["phy.platform"] != "physical":
            log.info("Not a physical platform, skipping bmc parsing")
            return

        # TODO: put in parse_fail_msg if the parse fails
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
            return

        for regex, field in regexes.items():
            match = re.search(regex, data, re.MULTILINE)
            value = match.group(1).strip() if match else None
            log.debug(f"{value=}")
            if value is not None:
                fields[field] = value
        log.debug(f"{fields=}")

        facts.update(fields)

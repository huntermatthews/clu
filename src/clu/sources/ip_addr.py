import json
import logging

from clu import Facts, Provides, Requires, Source
from clu.input import text_program

log = logging.getLogger(__name__)

primary_keys = ["net.macs", "net.ipv4", "net.ipv6", "net.devs"]


class IpAddr(Source):
    def provides(self, provides: Provides) -> None:
        for key in primary_keys:
            provides[key] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("ip --json addr")

    def parse(self, facts: Facts) -> None:
        output, rc = text_program("ip --json addr")
        log.debug(f"{output=}")
        if output == "" or rc != 0:
            for key in primary_keys:
                facts[key] = "Error/Unknown"
            return
        else:
            for key in primary_keys:
                facts[key] = ""

        data = json.loads(output)
        for iface in data:
            facts["net.devs"] += iface["ifname"] + " "
            for addr in iface["addr_info"]:
                if addr["family"] == "inet":
                    facts["net.ipv4"] += addr["local"] + " "
                elif addr["family"] == "inet6":
                    facts["net.ipv6"] += addr["local"] + " "
            facts["net.macs"] += iface["address"] + " "

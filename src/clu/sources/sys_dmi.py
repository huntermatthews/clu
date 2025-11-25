import logging

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts
from clu.sources import Source
from clu.input import text_file

log = logging.getLogger(__name__)


class SysDmi(Source):
    _dmap = {
        "sys.vendor": "/sys/devices/virtual/dmi/id/sys_vendor",
        "sys.model.family": "/sys/devices/virtual/dmi/id/product_family",
        "sys.model.name": "/sys/devices/virtual/dmi/id/product_name",
        "sys.serial_no": "/sys/devices/virtual/dmi/id/product_serial",
        "sys.uuid": "/sys/devices/virtual/dmi/id/product_uuid",
        "sys.oem": "/sys/devices/virtual/dmi/id/chassis_vendor",
        "sys.asset_no": "/sys/devices/virtual/dmi/id/chassis_asset_tag",
    }

    def provides(self, provides: Provides) -> None:
        for key in self._dmap.keys():
            provides[key] = self

    def requires(self, requires: Requires) -> None:
        requires.files.extend(self._dmap.values())

    def parse(self, facts: Facts) -> None:
        if facts["phy.platform"] != "physical":
            log.info("Not a physical platform, skipping sys.dmi parsing")
            return

        for key, entry in self._dmap.items():
            data = text_file(entry)
            log.debug(f"{data=}")
            log.debug(f"{key=}")
            facts[key] = data.strip() if data else ""

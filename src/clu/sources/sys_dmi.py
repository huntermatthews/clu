import logging

from clu import Facts, Provides, Requires, Source
from clu.sources.virt_what import VirtWhat
from clu.debug import panic
from clu.input import text_file

log = logging.getLogger(__name__)


class SysDmi(Source):
    def provides(self) -> Provides:
        provides = Provides()
        provides["sys.vendor"] = self
        provides["sys.model.family"] = self
        provides["sys.model.name"] = self
        provides["sys.serial_no"] = self
        provides["sys.uuid"] = self
        provides["sys.oem"] = self
        provides["sys.asset_no"] = self
        return provides

    def requires(self) -> Requires:
        requires = Requires()
        requires.files.extend(
            [
                "/sys/devices/virtual/dmi/id/sys_vendor",
                "/sys/devices/virtual/dmi/id/product_family",
                "/sys/devices/virtual/dmi/id/product_name",
                "/sys/devices/virtual/dmi/id/product_serial",
                "/sys/devices/virtual/dmi/id/product_uuid",
                "/sys/devices/virtual/dmi/id/chassis_vendor",
                "/sys/devices/virtual/dmi/id/chassis_asset_tag",
            ]
        )
        return requires

    def parse(self, facts: Facts) -> Facts:
        if "phy.platform" not in facts:
            virtwhat = VirtWhat()
            virtwhat.parse(facts)

        if facts["phy.platform"] != "physical":
            log.info("Not a physical platform, skipping sys.dmi parsing")
            return facts

        keys = [
            "sys.vendor",
            "sys.model.family",
            "sys.model.name",
            "sys.serial_no",
            "sys.uuid",
            "sys.oem",
            "sys.asset_no",
        ]
        entries = [
            "sys_vendor",
            "product_family",
            "product_name",
            "product_serial",
            "product_uuid",
            "chassis_vendor",
            "chassis_asset_tag",
        ]
        if len(keys) != len(entries):
            log.debug(f"{keys=}")
            log.debug(f"{entries=}")
            panic("parse_sys_dmi: keys and entries length don't match: You can't count")
        for idx in range(len(keys)):
            key = keys[idx]
            entry = entries[idx]
            data = text_file(f"/sys/devices/virtual/dmi/id/{entry}")
            log.debug(f"{data=}")
            log.debug(f"{key=}")
            facts[key] = data.strip() if data else ""

        return facts

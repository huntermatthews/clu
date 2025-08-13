import logging

from clu import Facts, Provides, Requires, panic
from clu.input import text_file

log = logging.getLogger(__name__)


def provides_sys_dmi(provides: Provides) -> None:
    provides["sys.vendor"] = parse_sys_dmi
    provides["sys.model.family"] = parse_sys_dmi
    provides["sys.model.name"] = parse_sys_dmi
    provides["sys.serial_no"] = parse_sys_dmi
    provides["sys.uuid"] = parse_sys_dmi
    provides["sys.oem"] = parse_sys_dmi
    provides["sys.asset_no"] = parse_sys_dmi


def requires_sys_dmi(requires: Requires) -> None:
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


def parse_sys_dmi(facts: Facts) -> None:
    if "phy.platform" not in facts:
        parse_virt_what(facts)

    if facts["phy.platform"] != "physical":
        log.info("Not a physical platform, skipping sys.dmi parsing")
        return

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

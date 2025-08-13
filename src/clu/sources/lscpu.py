import logging
import re

from clu import Facts, Provides, Requires
from clu.input import text_program

log = logging.getLogger(__name__)


def provides_lscpu(provides: Provides) -> None:
    provides["phy.cpu.model"] = parse_lscpu
    provides["phy.cpu.vendor"] = parse_lscpu
    provides["phy.cpu.cores"] = parse_lscpu
    provides["phy.cpu.threads"] = parse_lscpu
    provides["phy.cpu.sockets"] = parse_lscpu


def requires_lscpu(requires: Requires) -> None:
    requires.programs.append("lscpu")


# TODO: clean this up, it is a mess because it didn't translate from the original code well
def parse_lscpu(facts: Facts) -> None:
    regexes = {
        r"^ *Model name: *(.+)": "model",
        r"^ *Vendor ID: *(.+)": "vendor",
        r"^ *Core\(s\) per socket: *(\d+)": "cores_per_socket",
        r"^ *Thread\(s\) per core: *(\d+)": "threads_per_core",
        r"^ *Socket\(s\): *(\d+)": "sockets",
        r"^ *CPU\(s\): *(\d+)": "cpus",
    }
    fields = {}
    attr_keys = ["model", "vendor", "cores", "threads", "sockets"]
    data, rc = text_program("lscpu")
    if data is None or rc != 0:
        return

    for regex, field in regexes.items():
        match = re.search(regex, data, re.MULTILINE)
        value = match.group(1).strip() if match else None
        log.debug(f"{value=}")
        if value is not None:
            fields[field] = value
    log.debug(f"{fields=}")
    try:
        fields["cores"] = str(int(fields["cores_per_socket"]) * int(fields["sockets"]))
        fields["threads"] = str(int(fields["threads_per_core"]) * int(fields["cores"]))
    except Exception:
        pass
    log.debug(f"{fields=}")
    for key in attr_keys:
        value = fields.get(key)
        log.debug(f"{key=}")
        log.debug(f"{value=}")
        facts[f"phy.cpu.{key}"] = value

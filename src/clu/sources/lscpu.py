import logging
import re

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts
from clu.sources import Source, PARSE_FAIL_MSG
from clu.input import text_program

log = logging.getLogger(__name__)


class Lscpu(Source):
    def provides(self, provides: Provides) -> None:
        provides["phy.cpu.model"] = self
        provides["phy.cpu.vendor"] = self
        provides["phy.cpu.cores"] = self
        provides["phy.cpu.threads"] = self
        provides["phy.cpu.sockets"] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("lscpu")

    # TODO: clean this up, it is a mess because it didn't translate from the original code well
    def parse(self, facts: Facts) -> None:
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
        if data == "" or rc != 0:
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
            fields["cores"] = PARSE_FAIL_MSG
            fields["threads"] = PARSE_FAIL_MSG

        log.debug(f"{fields=}")
        for key in attr_keys:
            value = fields.get(key)
            log.debug(f"{key=}")
            log.debug(f"{value=}")
            facts[f"phy.cpu.{key}"] = value

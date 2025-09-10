import logging

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts
from clu.sources import Source
from clu.input import text_file
from clu.sources import PARSE_FAIL_MSG

log = logging.getLogger(__name__)


class OsRelease(Source):
    def provides(self, provides: Provides) -> None:
        provides["os.distro.name"] = self
        provides["os.distro.version"] = self

    def requires(self, requires: Requires) -> None:
        requires.files.append("/etc/os-release")

    def parse(self, facts: Facts) -> None:
        data = text_file("/etc/os-release")
        log.debug(f"{data=}")
        if not data:
            facts["os.distro.name"] = PARSE_FAIL_MSG
            facts["os.distro.version"] = PARSE_FAIL_MSG
            return

        for line in data.splitlines():
            if "=" not in line:
                continue
            key, value = line.split("=", 1)
            key = key.strip().strip('"')
            value = value.strip().strip('"')
            log.debug(f"{key=}")
            log.debug(f"{value=}")

            if key == "ID":
                facts["os.distro.name"] = value
            elif key == "VERSION_ID":
                facts["os.distro.version"] = value

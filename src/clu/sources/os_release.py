import logging

from clu import Facts, Provides, Requires, Source
from clu.input import text_file

log = logging.getLogger(__name__)


class OsRelease(Source):
    def provides(self, provides: Provides) -> None:
        provides["os.distro.name"] = self
        provides["os.distro.version"] = self

    def requires(self, requires: Requires) -> None:
        requires.files.append("/etc/os-release")

    def parse(self, facts: Facts) -> None:
        data = text_file("/etc/os-release")
        log.trace(f"{data=}")
        if not data:
            facts["os.distro.name"] = "Unknown/Error"
            facts["os.distro.version"] = "Unknown/Error"
            return
        
        for line in data.splitlines():
            if "=" not in line:
                continue
            key, value = line.split("=", 1)
            key = key.strip().strip('"')
            value = value.strip().strip('"')
            log.trace(f"{key=}")
            log.trace(f"{value=}")

            if key == "ID":
                facts["os.distro.name"] = value
            elif key == "VERSION_ID":
                facts["os.distro.version"] = value

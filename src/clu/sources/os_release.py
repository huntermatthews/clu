import logging

from clu import Facts, Provides, Requires, Source
from clu.input import text_file

log = logging.getLogger(__name__)


class OsRelease(Source):
    def provides(self) -> Provides:
        provides = Provides()
        provides["os.distro.name"] = self
        provides["os.distro.version"] = self
        return provides

    def requires(self) -> Requires:
        requires = Requires()
        requires.files.append("/etc/os-release")
        return requires

    def parse(self, facts: Facts) -> Facts:
        data = text_file("/etc/os-release")
        log.debug(f"{data=}")
        if not data:
            facts["os.distro.name"] = "Unknown/Error"
            facts["os.distro.version"] = "Unknown/Error"
            return facts
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

        return facts

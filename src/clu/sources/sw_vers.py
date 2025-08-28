import logging

from clu import Facts, Provides, Requires, Source
from clu.input import text_program

log = logging.getLogger(__name__)


class SwVers(Source):
    _keys = [
        "os.name",
        "os.version",
        "os.build",
    ]

    def provides(self, provides: Provides) -> None:
        for key in self._keys:
            provides[key] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("sw_vers")

    def parse(self, facts: Facts) -> None:
        data, rc = text_program("sw_vers")
        log.debug(f"{data=}")
        if data == "" or rc != 0:
            for key in self._keys:
                facts[key] = "Unknown/Error"
            return
        for line in data.splitlines():
            if ":" not in line:
                continue
            key, value = line.split(":", 1)
            key = key.strip()
            value = value.strip()
            log.debug(f"{key=}")
            log.debug(f"{value=}")

            if key == "ProductName":
                facts["os.name"] = value
            elif key == "ProductVersion":
                facts["os.version"] = value
            elif key == "BuildVersion":
                facts["os.build"] = value

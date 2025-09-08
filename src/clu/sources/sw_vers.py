import logging

from clu import facts, Tier, Provides, Requires
from clu.input import text_program
from clu.sources import Source, PARSE_FAIL_MSG

log = logging.getLogger(__name__)


class SwVers(Source):
    _names = {
        "ProductName": ("os.name", Tier.PRIMARY),
        "ProductVersion": ("os.version", Tier.PRIMARY),
        "BuildVersion": ("os.build", Tier.OTHER),
    }

    def provides(self, provides: Provides) -> None:
        for key in self._names:
            provides[self._names[key][0]] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("sw_vers")

    def parse(self) -> None:
        if "os.name" in facts:
            return

        data, rc = text_program("sw_vers")
        log.debug(f"{data=}")
        if data == "" or rc != 0:
            for name in self._names:
                facts.add(self._names[name][1], self._names[name][0], PARSE_FAIL_MSG)
            return

        for line in data.splitlines():
            if ":" not in line:
                continue
            key, value = line.split(":", 1)
            key = key.strip()
            value = value.strip()
            log.debug(f"{key=}")
            log.debug(f"{value=}")

            if key in self._names:
                facts.add(self._names[key][1], self._names[key][0], value)
            else:
                log.error(f"Unknown key: {key}")

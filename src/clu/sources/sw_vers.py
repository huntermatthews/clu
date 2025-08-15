import logging

from clu import Facts, Provides, Requires, Source
from clu.input import text_program

log = logging.getLogger(__name__)


class SwVers(Source):
    def provides(self) -> Provides:
        provides = Provides()

        for key in [
            "os.name",
            "os.version",
            "os.build",
        ]:
            provides[key] = self

        return provides

    def requires(self) -> Requires:
        requires = Requires()

        requires.programs.append("sw_vers")

        return requires

    def parse(self, facts: Facts) -> Facts:
        if "os.version" in facts:
            return facts

        data, rc = text_program("sw_vers")
        log.debug(f"{data=}")
        if data is None or rc != 0:
            return facts
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

        return facts

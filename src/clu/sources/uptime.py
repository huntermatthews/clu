import logging
import re

from clu import Facts, Provides, Requires, Source
from clu.input import text_program

log = logging.getLogger(__name__)


class Uptime(Source):
    def provides(self, provides: Provides) -> None:
        provides["run.uptime"] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("uptime")

    def parse(self, facts: Facts) -> None:
        data, rc = text_program("uptime")
        log.debug(f"{data=}")
        if data == "" or rc != 0:
            facts["run.uptime"] = "Unknown/Error"
            return

        match = re.match(r".*up *(.*) \d+ user.*", data)
        uptime = match.group(1).rstrip(",") if match else "Error/Unknown"
        log.debug(f"{uptime=}")
        facts["run.uptime"] = uptime

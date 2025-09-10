import logging
import re

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts
from clu.sources import Source
from clu.input import text_program
from clu.sources import PARSE_FAIL_MSG

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
            facts["run.uptime"] = PARSE_FAIL_MSG
            return

        match = re.match(r".*up *(.*) \d+ user.*", data)
        uptime = match.group(1).rstrip(",") if match else PARSE_FAIL_MSG
        log.debug(f"{uptime=}")
        facts["run.uptime"] = uptime

import logging
import re

from clu import facts, Tier, Provides, Requires
from clu.input import text_program
from clu.sources import Source, PARSE_FAIL_MSG

log = logging.getLogger(__name__)


class Uptime(Source):
    def provides(self, provides: Provides) -> None:
        provides["run.uptime"] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.append("uptime")

    def parse(self) -> None:
        data, rc = text_program("uptime")
        log.debug(f"{data=}")
        if data == "" or rc != 0:
            facts.add(Tier.PRIMARY, "run.uptime", PARSE_FAIL_MSG)
            return

        match = re.match(r".*up *(.*) \d+ user.*", data)
        uptime = match.group(1).rstrip(",") if match else PARSE_FAIL_MSG
        log.debug(f"{uptime=}")
        facts.add(Tier.PRIMARY, "run.uptime", uptime)

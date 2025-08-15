import logging
import re

from clu import Facts, Provides, Requires, Source
from clu.debug import panic
from clu.input import text_program

log = logging.getLogger(__name__)


class Uptime(Source):
    def provides(self) -> Provides:
        provides = Provides()
        provides["run.uptime"] = self.parse
        return provides

    def requires(self) -> Requires:
        requires = Requires()
        requires.programs.append("uptime")
        return requires

    def parse(self) -> Facts:
        facts = Facts()
        data, rc = text_program("uptime")
        if data is None or rc != 0:
            panic("parse_uptime: uptime command failed")
        log.debug(f"{data=}")
        match = re.match(r".*up *(.*) \d+ user.*", data)
        uptime = match.group(1).rstrip(",") if match else "unknown / error"
        log.debug(f"{uptime=}")
        facts["run.uptime"] = uptime
        return facts

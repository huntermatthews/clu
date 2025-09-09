import logging

from clu import facts, Provides, Requires
from clu.input import text_file
from clu.conversions import seconds_to_text
from clu.sources import Source, PARSE_FAIL_MSG

log = logging.getLogger(__name__)


class ProcUptime(Source):
    def provides(self, provides: Provides) -> None:
        provides["run.uptime"] = self

    def requires(self, requires: Requires) -> None:
        requires.files.append("/proc/uptime")

    def parse(self) -> None:
        data = text_file("/proc/uptime")
        log.debug(f"{data=}")
        if not data:
            facts["run.uptime"] = PARSE_FAIL_MSG
            return
        uptime_secs = int(float(data.split()[0]))
        log.debug(f"{uptime_secs=}")
        facts["run.uptime"] = seconds_to_text(uptime_secs)

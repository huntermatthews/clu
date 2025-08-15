import logging

from clu import Facts, Provides, Requires, Source
from clu.input import text_file
from clu.conversions import seconds_to_text

log = logging.getLogger(__name__)


class ProcUptime(Source):
    def provides(self) -> Provides:
        provides = Provides()
        provides["run.uptime"] = self
        return provides

    def requires(self) -> Requires:
        requires = Requires()
        requires.files.append("/proc/uptime")
        return requires

    def parse(self, facts: Facts) -> Facts:
        data = text_file("/proc/uptime")
        log.debug(f"{data=}")
        if not data:
            return facts
        uptime_secs = int(float(data.split()[0]))
        log.debug(f"{uptime_secs=}")
        facts["run.uptime"] = seconds_to_text(uptime_secs)
        return facts

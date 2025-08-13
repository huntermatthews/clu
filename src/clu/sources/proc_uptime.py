import logging

from clu import Facts, Provides, Requires
from clu.input import text_file
from clu.conversions import seconds_to_text

log = logging.getLogger(__name__)


def provides_proc_uptime(provides: Provides) -> None:
    provides["run.uptime"] = parse_proc_uptime


def requires_proc_uptime(requires: Requires) -> None:
    requires.files.append("/proc/uptime")


def parse_proc_uptime(facts: Facts) -> None:
    data = text_file("/proc/uptime")
    log.debug(f"{data=}")
    if not data:
        return
    uptime_secs = int(float(data.split()[0]))
    log.debug(f"{uptime_secs=}")
    facts["run.uptime"] = seconds_to_text(uptime_secs)

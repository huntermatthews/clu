import logging

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts
from clu.sources import Source
from clu.input import text_program
from clu.sources import PARSE_FAIL_MSG, NET_DISABLED_MSG
from clu.config import get_config

log = logging.getLogger(__name__)
cfg = get_config()


class DnfCheckUpdate(Source):
    def provides(self, provides: Provides) -> None:
        provides["run.update_required"] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.extend(["dnf check-update"])

    def parse(self, facts: Facts) -> None:
        if not cfg.net:
            facts["run.update_required"] = NET_DISABLED_MSG
            return

        _, rc = text_program("dnf check-update")
        log.debug(f"rc is {rc}")

        if rc == 0:
            value = "False"
        elif rc == 100:
            value = "True"
        else:
            value = PARSE_FAIL_MSG

        facts["run.update_required"] = value

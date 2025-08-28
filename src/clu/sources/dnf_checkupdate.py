import logging

from clu import Facts, Provides, Requires, Source
from clu.input import text_program

log = logging.getLogger(__name__)


class DnfCheckUpdate(Source):
    def provides(self, provides: Provides) -> None:
        provides["run.update_required"] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.extend(["dnf check-update"])

    def parse(self, facts: Facts) -> None:
        _, rc = text_program("dnf check-update")
        log.debug(f"rc is {rc}")
        if rc == 0:
            facts["run.update_required"] = "False"
        elif rc == 100:
            facts["run.update_required"] = "True"
        else:
            facts["run.update_required"] = "Unknown/Error"

import logging

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts
from clu.sources import Source
from clu.input import text_file

log = logging.getLogger(__name__)


class NoSalt(Source):
    def provides(self, provides: Provides) -> None:
        provides["salt.no_salt.exists"] = self
        provides["salt.no_salt.reason"] = self

    def requires(self, requires: Requires) -> None:
        requires.files.append("/no_salt")

    def parse(self, facts: Facts) -> None:
        data = text_file("/no_salt", optional=True)
        log.debug(f"{data=}")
        if data == "":
            facts["salt.no_salt.exists"] = "False"
            return
        else:
            log.debug(f"{data=}")
            facts["salt.no_salt.exists"] = "True"

        facts["salt.no_salt.reason"] = data.strip()

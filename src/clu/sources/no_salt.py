import logging

from clu import Facts, Provides, Requires, Source
from clu.input import text_file

log = logging.getLogger(__name__)


class NoSalt(Source):
    def provides(self) -> Provides:
        provides = Provides()
        provides["os.no_salt.exists"] = self
        provides["os.no_salt.reason"] = self
        return provides

    def requires(self) -> Requires:
        requires = Requires()
        requires.files.append("/no_salt")
        return requires

    def parse(self, facts: Facts) -> Facts:
        data = text_file("/no_salt")
        if data is None:
            facts["salt.no_salt.exists"] = "False"
            return facts
        else:
            log.debug(f"{data=}")
            facts["salt.no_salt.exists"] = "True"
        if not data.strip():
            data = "UNKNOWN"
        facts["salt.no_salt.reason"] = data.strip()

        return facts

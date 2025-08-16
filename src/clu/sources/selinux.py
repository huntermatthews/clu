import logging

from clu import Facts, Provides, Requires, Source
from clu.input import text_program

log = logging.getLogger(__name__)


class Selinux(Source):
    def provides(self, provides: Provides) -> None:
        provides["os.selinux.enable"] = self
        provides["os.selinux.mode"] = self

    def requires(self, requires: Requires) -> None:
        requires.programs.extend(["selinuxenabled", "getenforce"])

    def parse(self, facts: Facts) -> None:
        _, rc = text_program("selinuxenabled")
        # man page: "status 0 if SELinux is enabled and 1 if it is not enabled."
        log.debug(f"rc is {rc}")
        if rc == 0:
            facts["os.selinux.enable"] = "True"
        elif rc == 1:
            facts["os.selinux.enable"] = "False"
        else:
            facts["os.selinux.enable"] = "Unknown/Error"

        data, rc = text_program("getenforce")
        log.debug(f"{data=}")
        facts["os.selinux.mode"] = data.strip() if data else "Unknown/Error"

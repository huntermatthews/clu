import logging

from clu import Facts, Provides, Requires, Source
from clu.input import text_program

log = logging.getLogger(__name__)


class Selinux(Source):
    def provides(self) -> Provides:
        provides = Provides()
        provides["os.selinux.enable"] = self
        provides["os.selinux.mode"] = self
        return provides

    def requires(self) -> Requires:
        requires = Requires()
        requires.programs.extend(["selinuxenabled", "getenforce"])
        return requires

    def parse(self, facts: Facts) -> Facts:
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
        return facts

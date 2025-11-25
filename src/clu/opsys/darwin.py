import logging

from clu.opsys import OpSys
from clu.sources import uname, system_version_plist, macos_name, uptime, clu

log = logging.getLogger(__name__)


class Darwin(OpSys):
    """macOS (Darwin) operating system class."""

    _sources = [
        uname.Uname(),
        system_version_plist.SystemVersionPlist(),
        macos_name.MacOSName(),
        uptime.Uptime(),
        clu.Clu(),
    ]

    def default_facts(self) -> list:
        return ["os.name", "os.hostname", "os.version", "os.code_name", "run.uptime", "clu.version"]

    def early_facts(self) -> list:
        return ["os.version"]

    # def parse(self, facts):
    #     """Parse the facts for macOS (Darwin)."""
    #     facts = super().parse(facts)
    #     # Nothing explicitly says Apple, but we know its apple because Darwin is the OS
    #     facts["sys.vendor"] = "Apple"
    #     return facts

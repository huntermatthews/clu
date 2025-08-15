"""Doc Incomplete."""

import logging

from clu import Facts
from clu.opsys import OpSys

from clu.sources import uname, sw_vers, macos_name, uptime, clu

log = logging.getLogger(__name__)


class Darwin(OpSys):
    """macOS (Darwin) operating system class."""

    _sources = [uname.Uname, sw_vers.SwVers, macos_name.MacOSName, uptime.Uptime, clu.Clu]

    def default_facts(self) -> list:
        return ["os.name", "os.hostname", "os.version", "os.codename", "run.uptime", "clu.version"]

    def parse(self) -> Facts:
        """Parse the facts for macOS (Darwin)."""
        facts = super().parse()

        # Nothing explicitly says Apple, but we know its apple because Darwin is the OS
        facts["sys.vendor"] = "Apple"

        return facts

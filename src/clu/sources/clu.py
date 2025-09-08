import datetime
import getpass
import logging
import os
import sys

from clu import facts, Tier, Provides, Requires, __about__
from clu.sources import Source

log = logging.getLogger(__name__)


class Clu(Source):
    def provides(self, provides: Provides) -> None:
        """Define the provider map for CLU."""
        for key in [
            "clu.binary",
            "clu.version",
            "clu.python.binary",
            "clu.python.version",
            "clu.cmdline",
            "clu.cwd",
            "clu.user",
            "clu.date",
        ]:
            provides[key] = self

    def requires(self, requires: Requires) -> None:
        # No specific requirements for clu group
        pass

    def _get_rfc3339_timestamp(self) -> str:
        """Get the current time in RFC 3339 format."""
        now_utc = datetime.datetime.now(datetime.timezone.utc)
        return now_utc.isoformat(sep="T", timespec="seconds")

    def parse(self) -> None:
        """Return info about CLU itself (mostly runtime)."""
        facts.add(Tier.SECONDARY, "clu.binary", sys.argv[0])
        facts.add(Tier.PRIMARY, "clu.version", __about__.__version__)
        facts.add(Tier.SECONDARY, "clu.python.binary", sys.executable)
        facts.add(Tier.PRIMARY, "clu.python.version", ".".join(map(str, sys.version_info[:3])))
        facts.add(Tier.SECONDARY, "clu.cmdline", " ".join(sys.argv))
        facts.add(Tier.SECONDARY, "clu.cwd", os.getcwd())
        facts.add(Tier.SECONDARY, "clu.user", getpass.getuser())
        facts.add(Tier.SECONDARY, "clu.date", self._get_rfc3339_timestamp())

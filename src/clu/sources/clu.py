import datetime
import getpass
import logging
import os
import sys

from clu import __about__
from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts, Tier
from clu.sources import Source

log = logging.getLogger(__name__)


def get_rfc3339_timestamp() -> str:
    """Get the current time in RFC 3339 format."""
    now_utc = datetime.datetime.now(datetime.timezone.utc)
    return now_utc.isoformat(sep="T", timespec="seconds")


# This silly function makes testing easier - its easier to patch() out
def raw_clu_metadata() -> dict:
    """Return info about CLU itself (mostly runtime)."""

    return {
        "sys.argv": sys.argv,
        "about.version": __about__.__version__,
        "sys.executable": sys.executable,
        "sys.version_info": sys.version_info,
        "os.getcwd": os.getcwd(),
        "getpass.getuser": getpass.getuser(),
        "datestamp": get_rfc3339_timestamp(),
    }


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

    def parse(self, facts: Facts) -> None:
        """Return info about CLU itself (mostly runtime)."""

        metadata = raw_clu_metadata()
        facts.add(Tier.TWO, "clu.binary", metadata["sys.argv"][0])
        facts.add(Tier.ONE, "clu.version", metadata["about.version"])
        facts.add(Tier.TWO, "clu.python.binary", metadata["sys.executable"])
        facts.add(
            Tier.ONE, "clu.python.version", ".".join(map(str, metadata["sys.version_info"][:3]))
        )
        facts.add(Tier.TWO, "clu.cmdline", " ".join(metadata["sys.argv"]))
        facts.add(Tier.THREE, "clu.cwd", metadata["os.getcwd"])
        facts.add(Tier.THREE, "clu.user", metadata["getpass.getuser"])
        facts.add(Tier.TWO, "clu.date", metadata["datestamp"])

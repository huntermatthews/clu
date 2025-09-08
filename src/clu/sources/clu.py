import datetime
import getpass
import logging
import os
import sys

from clu import Facts, Provides, Requires, Source, __about__

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

    def parse(self, facts: Facts) -> None:
        """Return info about CLU itself (mostly runtime)."""
        facts["clu.binary"] = sys.argv[0]
        facts["clu.version"] = __about__.__version__
        facts["clu.python.binary"] = sys.executable
        facts["clu.python.version"] = ".".join(map(str, sys.version_info[:3]))
        facts["clu.cmdline"] = " ".join(sys.argv)
        facts["clu.cwd"] = os.getcwd()
        facts["clu.user"] = getpass.getuser()
        facts["clu.date"] = self._get_rfc3339_timestamp()

        facts["_primary"].extend(
            [
                "clu.version",
                "clu.python.version",
            ]
        )

        facts["_secondary"].extend(
            [
                "clu.binary",
                "clu.python.binary",
                "clu.cmdline",
                "clu.cwd",
                "clu.user",
                "clu.date",
            ]
        )

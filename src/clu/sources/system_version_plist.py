import logging
import plistlib

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts, Tier
from clu.sources import Source, PARSE_FAIL_MSG
from clu.input import text_file

log = logging.getLogger(__name__)


class SystemVersionPlist(Source):
    _keys = [
        "os.name",
        "os.version",
        "os.build",
        "id.build_id",
    ]

    def provides(self, provides: Provides) -> None:
        for key in self._keys:
            provides[key] = self

    def requires(self, requires: Requires) -> None:
        requires.files.append("/System/Library/CoreServices/SystemVersion.plist")

    def parse(self, facts: Facts) -> None:
        if "os.name" in facts:
            return

        data = text_file("/System/Library/CoreServices/SystemVersion.plist")
        plist = plistlib.loads(bytes(data, "utf-8"))
        log.debug(f"{plist=}")

        if not plist:
            for key in self._keys:
                facts[key] = PARSE_FAIL_MSG
            return
        else:
            facts.add(Tier.ONE, "os.name", plist.get("ProductName", PARSE_FAIL_MSG))
            facts.add(Tier.ONE, "os.version", plist.get("ProductVersion", PARSE_FAIL_MSG))
            facts.add(Tier.TWO, "os.build", plist.get("ProductBuildVersion", PARSE_FAIL_MSG))
            facts.add(Tier.THREE, "id.build_id", plist.get("BuildID", PARSE_FAIL_MSG))

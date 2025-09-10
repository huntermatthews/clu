import logging

from clu.provides import Provides
from clu.requires import Requires
from clu.sources import Source

log = logging.getLogger(__name__)


class OpSys:
    _sources: list[Source] = []

    def default_facts(self) -> list[str]:
        return []

    def early_facts(self) -> list[str]:
        return []

    def provides(self) -> Provides:
        """Define the provider map for macOS (Darwin)."""
        provs = Provides()

        for source in self._sources:
            log.info(f"Source: {source}")
            log.debug(f"Source type: {type(source)}")
            source.provides(provs)

        return provs

    def requires(self) -> Requires:
        """Define the requirements for macOS (Darwin)."""
        reqs = Requires()

        for source in self._sources:
            source.requires(reqs)

        return reqs

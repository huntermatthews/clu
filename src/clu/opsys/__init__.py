from clu import Provides, Requires
from clu.debug import panic


class OpSys:
    _sources: list = []

    def default_facts(self) -> list[str]:
        return []

    def provides(self) -> Provides:
        """Define the provider map for macOS (Darwin)."""
        provs = Provides()

        for source in self._sources:
            print(source)
            provs.update(source.provides())

        return provs

    def requires(self) -> Requires:
        """Define the requirements for macOS (Darwin)."""
        reqs = Requires()

        for source in self._sources:
            reqs.update(source.requires())

        return reqs

    def parse(self):
        panic("Should not call parse on OpSys objects")

    # def parse(self) -> Facts:
    #     """Parse the facts for macOS (Darwin)."""
    #     facts = Facts()

    #     for source in self._sources:
    #         facts.update(source.parse())

    #     return facts

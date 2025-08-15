from clu import Provides, Requires, Facts, Source


class OpSys:
    _sources: list = []

    def default_facts(self) -> list[str]:
        return []

    def provides(self) -> Provides:
        """Define the provider map for macOS (Darwin)."""
        provides = Provides()

        for source in self._sources:
            provides.update(source.provides())

        return provides

    def requires(self) -> Requires:
        """Define the requirements for macOS (Darwin)."""
        requires = Requires()

        for source in self._sources:
            requires.update(source.requires())

        return requires

    def parse(self) -> Facts:
        """Parse the facts for macOS (Darwin)."""
        facts = Facts()

        for source in self._sources:
            facts.update(source.parse())

        return facts

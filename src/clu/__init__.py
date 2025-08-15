from dataclasses import dataclass, field
from typing import List


class Facts(dict):
    """A simple facts container that behaves like a dict."""

    pass


class Provides(dict):
    """A simple provides container that behaves like a dict."""

    pass


@dataclass
class Requires:
    files: List[str] = field(default_factory=list)
    programs: List[str] = field(default_factory=list)
    apis: List[str] = field(default_factory=list)
    facts: List[str] = field(default_factory=list)

    def update(self, other: "Requires"):
        self.files.extend(other.files)
        self.programs.extend(other.programs)
        self.apis.extend(other.apis)
        self.facts.extend(other.facts)
        return self

    def sort(self):
        self.files.sort()
        self.programs.sort()
        self.apis.sort()
        self.facts.sort()
        return self


class Source:
    def provides(self) -> Provides:
        """Define the provider map for this source."""
        return Provides()

    def requires(self) -> Requires:
        """Define the requirements for this source."""
        return Requires()

    def parse(self, facts: Facts) -> Facts:
        """Parse the facts for this source."""
        print("im being parsed!", __class__)
        return facts

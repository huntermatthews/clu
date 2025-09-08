from dataclasses import dataclass, field
from typing import List

from clu.facts import Facts, Tier


facts: Facts = Facts()


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


__all__ = [
    "Facts",
    "facts",
    "Provides",
    "Requires",
    "Tier",
]

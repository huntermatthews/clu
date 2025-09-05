from dataclasses import dataclass, field
from abc import ABC, abstractmethod
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

    def update(self, other: "Requires"):
        self.files.extend(other.files)
        self.programs.extend(other.programs)
        self.apis.extend(other.apis)
        return self


class Source(ABC):
    @abstractmethod
    def provides(self, provides: Provides):
        """Define the provider map for this source."""
        pass

    @abstractmethod
    def requires(self, requires: Requires):
        """Define the requirements for this source."""
        pass

    @abstractmethod
    def parse(self, facts: Facts):
        """Parse the facts for this source."""
        pass

from abc import ABC, abstractmethod

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts


PARSE_FAIL_MSG = "Unknown/Error"
NET_DISABLED_MSG = "Unknown - Network Queries Disabled"


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

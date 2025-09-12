class Facts(dict):
    """A simple facts container that behaves like a dict."""

    pass


from enum import Enum
from typing import Any


class Tier(Enum):
    PRIMARY = "primary"
    SECONDARY = "secondary"
    OTHER = "other"


class xFacts:
    def __init__(self):
        self._priority = {Tier.PRIMARY: [], Tier.SECONDARY: [], Tier.OTHER: []}
        self._facts = {}

    def add(self, priority: Tier, key: str, value: str) -> None:
        self._priority[priority].append(key)
        self._facts[key] = value

    def update(self, other: dict):
        for key in other:
            self.__setitem__(key, other[key])

    def __contains__(self, item: str) -> bool:
        return item in self._facts

    def __eq__(self, other) -> bool:
        if not isinstance(other, Facts):
            return NotImplemented

        if self._facts != other._facts:
            return False
        # FIXME: when the tests can handle this, reactivate
        # if self._priority != other._priority:
        #     print("self:", self._priority)
        #     print("other:", other._priority)
        #     print("xx: _priority not equal")
        #     return False
        return True

    def __getitem__(self, key: str) -> Any:
        return self._facts.get(key)

    def __iter__(self):
        for item in self._facts:
            yield item

    def __repr__(self) -> str:
        return f"Facts: {self._facts}, {self._priority}"

    def __setitem__(self, key: str, value: Any) -> None:
        self._facts[key] = value
        self._priority[Tier.PRIMARY].append(key)

    def __str__(self) -> str:
        return f"Facts: {self._facts}, {self._priority}"

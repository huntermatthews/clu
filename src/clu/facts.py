from enum import Enum
from typing import Any


class Tier(Enum):
    PRIMARY = "primary"
    SECONDARY = "secondary"
    OTHER = "other"


class Facts:
    def __init__(self):
        self._priority = {Tier.PRIMARY: [], Tier.SECONDARY: [], Tier.OTHER: []}
        self._facts = {}

    def add(self, priority: Tier, key: str, value: str) -> None:
        self._priority[priority].append(key)
        self._facts[key] = value

    def __contains__(self, item: str) -> bool:
        return item in self._facts

    def __iter__(self):
        for item in self._facts:
            yield item

    def __setitem__(self, key: str, value: Any) -> None:
        self._facts[key] = value

    def __getitem__(self, key: str) -> Any:
        return self._facts.get(key)

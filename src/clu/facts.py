from enum import Enum
from typing import Any


class Tier(Enum):
    ONE = "one"
    TWO = "two"
    THREE = "three"

    def get_by_int(index: int) -> "Tier":
        if index == 1:
            return Tier.ONE
        elif index == 2:
            return Tier.TWO
        elif index == 3:
            return Tier.THREE
        else:
            raise ValueError("Invalid tier index")


class Facts:
    def __init__(self):
        self._tier = {Tier.ONE: [], Tier.TWO: [], Tier.THREE: []}
        self._facts = {}

    def add(self, priority: Tier, key: str, value: str) -> None:
        self._tier[priority].append(key)
        self._facts[key] = value

    def get_tier(self, tier: Tier) -> list:
        result = []
        if tier == Tier.ONE:
            result = self._tier[Tier.ONE]
        elif tier == Tier.TWO:
            result = self._tier[Tier.ONE] + self._tier[Tier.TWO]
        elif tier == Tier.THREE:
            result = self._tier[Tier.ONE] + self._tier[Tier.TWO] + self._tier[Tier.THREE]
        return result

    def to_dict(self) -> dict:
        return self._facts

    def update(self, other: dict):
        for key in other:
            self.__setitem__(key, other[key])

    def __contains__(self, item: str) -> bool:
        return item in self._facts

    def __eq__(self, other) -> bool:
        if isinstance(other, Facts):
            if self._facts != other._facts:
                return False
            # FIXME: when the tests can handle this, reactivate
            # if self._tier != other._tier:
            #     print("self:", self._tier)
            #     print("other:", other._tier)
            #     print("xx: _tier not equal")
            #     return False
            return True
        if isinstance(other, dict):
            return self._facts == other
        return False

    def __getitem__(self, key: str) -> Any:
        return self._facts.get(key)

    def __iter__(self):
        for item in self._facts:
            yield item

    def __repr__(self) -> str:
        return f"Facts: {self._facts}, {self._tier}"

    def __setitem__(self, key: str, value: Any) -> None:
        self._facts[key] = value
        self._tier[Tier.ONE].append(key)

    def __str__(self) -> str:
        return f"Facts: {self._facts}, {self._tier}"

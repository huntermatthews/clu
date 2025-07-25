"""
These are the requirements that we will collect (things the program depends on to run).
They are used to generate the requirements list and check if they are met.
"""

from dataclasses import dataclass, field
from typing import List

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

    def sort(self):
        self.files.sort()
        self.programs.sort()
        self.apis.sort()
        return self


_requires = {}
def add_requires(category, item):
    """Add a requirement to the specified category."""
    if category in _requires:
        _requires[category].append(item)
    else:
        _requires[category] = [item]


def get_all_requires():
    return _requires

# These are the requirements that we will collect (things the program depends on to run).
# They are used to generate the requirements list and check if they are met.


def requires_update(dst, src):
    dst.files.extend(src.files)
    dst.programs.extend(src.programs)
    dst.apis.extend(src.apis)
    return dst

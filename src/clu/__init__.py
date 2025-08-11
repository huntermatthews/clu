
"""Doc Incomplete."""

import argparse
import logging
from dataclasses import dataclass, field
from typing import List


log = logging.getLogger(__name__)

config = argparse.Namespace()


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


class Facts(dict):
	"""A simple facts container that behaves like a dict."""
	pass


class Provides(dict):
	"""A simple provides container that behaves like a dict."""
	pass


"""This file contains metadata about the package."""

from importlib import metadata

##
## Static fields
##

# Recall that __name__ is a "keyword" and thus can't be used here.
__title__ = "clu"
__version__ = "0.1.0"  # This should be updated manually or via a build script.
__url__ = "https://github.com/hunterjmatthews/clu"

__author__ = "Hunter Matthews"
__email__ = "hunter@unix.haus"

__copyright__ = "(c) 2012-2025 Hunter Matthews"
__license__ = "Apache-2.0"


##
## Computed fields
##
# try:
#     __version__ = metadata.version(__title__) or "0.dev0+unknown"
# except Exception:
#     __version__ = "0.dev0+unknown"

# # Our makefile/build process creates this file.
# try:
#     from __build_date__ import __build_date__
# except ImportError:
#     __build_date__ = "1970-01-01"

# # The rest of this we pull from the package's metadata (mostly to avoid duplication)
# dist_info = metadata.distribution(__title__)

# # NOTE: This mess is due to:
# #        a) you  can have non-direct dependencies and
# #        b) github based deps are 'name SPACE messy_url' - and we only want the name.
# __requirements__ = [
#     req.split(" ")[0] for req in dist_info.requires if "extra ==" not in req
# ]
# __summary__ = dist_info.metadata["Summary"]

# # Make some assumptions about the metadata for the requires-python field.
# assert dist_info.metadata["Requires-Python"].startswith(">=")
# __minimum_python__ = dist_info.metadata["Requires-Python"][2:]

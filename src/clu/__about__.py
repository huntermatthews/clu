"""This file contains metadata about the package."""

from importlib import metadata

__all__ = [
    "__title__",
    "__version__",
    "__email__",
    "__license__",
    "__summary__",
    "__minimum_python__",
    "__minimum_python_info__",
]

# Recall that __name__ is a "keyword" and thus can't be used here.
__title__: str = "clu"

try:
    __version__ = metadata.version(__title__) or "0.dev0+unknown"
except Exception:
    __version__ = "0.dev0+unknown"

# The rest of this we pull from the package's metadata (mostly to avoid duplication)
dist_info = metadata.distribution(__title__)

__summary__ = dist_info.metadata["Summary"]
__license__ = dist_info.metadata["License-Expression"]
__email__ = dist_info.metadata["Author-email"]

# Make some assumptions about the metadata for the requires-python field.
assert dist_info.metadata["Requires-Python"].startswith(">=")
__minimum_python__: str = dist_info.metadata["Requires-Python"][2:]
__minimum_python_info__: tuple = tuple([int(item) for item in __minimum_python__.split(".")])


# # NOTE: This mess is due to:
# #        a) you  can have non-direct dependencies and
# #        b) github based deps are 'name SPACE messy_url' - and we only want the name.
# __requirements__ = [
#     req.split(" ")[0] for req in dist_info.requires if "extra ==" not in req
# ]

# project_urls = dist_info.metadata.get_all('Project-URL')
# if project_urls:
#     for url_entry in project_urls:
#         label, url = url_entry.split(', ', 1) # Split the label and URL
#         print(f"Label: {label}, URL: {url}")
# else:
#     print("No Project-URL found for this package.")

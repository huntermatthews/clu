# These are the requirements that we will collect (things the program depends on to run).
# They are used to generate the requirements list and check if they are met.
_requires = {
    "files": [],
    "programs": [],
    "apis": [],
}

def add_requires(category, item):
    """Add a requirement to the specified category."""
    if category in _requires:
        _requires[category].append(item)
    else:
        _requires[category] = [item]


def get_all_requires():
    return _requires

"""Doc Incomplete."""

import argparse


# I can global if I want to
# (remember, this was ported from a couple of scripts that used globals)
facts = {}
config = argparse.Namespace()

# These are the requirements that we will collect
# They are used to generate the requirements list and check if they are met.
requires = {
    "files": [],
    "programs": [],
    "apis": [],
}

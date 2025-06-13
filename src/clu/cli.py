""" The main() function implementation for program.
"""

import argparse
import sys

from .__about__ import __version__


def main():
    if sys.version_info < (3,8):
        print("ERROR: Must use at least python 3.8")
        sys.exit(1)

    # TODO : check for root and bomb -- give up.

    parser = argparse.ArgumentParser(
        description="Show various information about the host",
        prog="clu",
        )
    parser.add_argument("--version", action="version", version="%%(prog)s %s" % __version__)
    parser.add_argument("facts", nargs="*")
    args = parser.parse_args()

    print(f"DEBUG: args = {args}")
    print(f"DEBUG: facts = {args.facts}")

    print("code here")


